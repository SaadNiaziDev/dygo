// Package migration owns the database transaction for project metadata, App
// installation, patches, access, and Fixtures. Local App commands never call it.
package migration

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hapyco/dygo/internal/access"
	"github.com/hapyco/dygo/internal/app/registry"
	"github.com/hapyco/dygo/internal/corevalues"
	"github.com/hapyco/dygo/internal/db"
	"github.com/hapyco/dygo/internal/fixtures"
	"github.com/hapyco/dygo/internal/patches"
	"github.com/hapyco/dygo/internal/project"
	"github.com/hapyco/dygo/internal/routes"
)

type AppChange struct{ Name, Version, Before, After string }
type Plan struct {
	PreSync        db.PatchPlan
	Schema         db.SchemaPlan
	SchemaDeferred bool
	PostSync       db.PatchPlan
	Access         access.Plan
	Fixtures       fixtures.Plan
	Apps           []AppChange
	fingerprint    [32]byte
	metadata       project.RuntimeMetadata
	installing     map[string]bool
	loadedPatches  []patches.LoadedPatch
	runs           []db.PatchRun
}
type Result struct {
	PreSync  db.PatchApplyResult
	Schema   db.SchemaSyncResult
	PostSync db.PatchApplyResult
	Access   access.Result
	Fixtures fixtures.Result
}
type Runner struct {
	Migrator    db.Migrator
	RecordHooks *db.RecordHookRegistry
	// ValidateCode verifies compiled registrations without generating or compiling code.
	ValidateCode func(string, project.RuntimeMetadata) error
}

func (r Runner) Plan(ctx context.Context, root, databaseURL string) (Plan, error) {
	pool, err := db.OpenRuntimePool(ctx, databaseURL)
	if err != nil {
		return Plan{}, err
	}
	defer pool.Close()
	return r.plan(ctx, pool, root)
}

func (r Runner) plan(ctx context.Context, q db.SchemaQueryer, root string) (Plan, error) {
	p := Plan{installing: map[string]bool{}}
	metadata, err := project.LoadRuntimeMetadata(root)
	if err != nil {
		return p, err
	}
	metadata.Apps, err = registry.DependencyOrder(metadata.Apps)
	if err != nil {
		return p, err
	}
	if _, err := routes.ValidateWithPages(metadata.Entities, metadata.Pages); err != nil {
		return p, err
	}
	if r.ValidateCode != nil {
		if err := r.ValidateCode(root, metadata); err != nil {
			return p, err
		}
	}
	p.metadata = metadata
	live, err := db.InspectLiveSchema(ctx, q)
	if err != nil {
		return p, err
	}
	states := map[string]AppChange{}
	if _, exists := live.Tables["app"]; exists {
		rows, err := q.Query(ctx, `SELECT name, version, status FROM "app" ORDER BY name`)
		if err != nil {
			return p, err
		}
		for rows.Next() {
			var a AppChange
			if err := rows.Scan(&a.Name, &a.Version, &a.Before); err != nil {
				rows.Close()
				return p, err
			}
			states[a.Name] = a
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return p, err
		}
	}
	found := map[string]bool{}
	for _, app := range metadata.Apps {
		found[app.Manifest.Name] = true
	}
	if !found["core"] {
		return p, fmt.Errorf("Core source is required for database migration")
	}
	for name := range states {
		if !found[name] {
			return p, fmt.Errorf("installed App %q source is missing; restore it before migration", name)
		}
	}
	for _, app := range metadata.Apps {
		name := app.Manifest.Name
		prior, exists := states[name]
		after := corevalues.AppStatusActive
		switch prior.Before {
		case "", corevalues.AppStatusInstalled:
			if exists && prior.Version != app.Manifest.Version {
				return p, fmt.Errorf("App %s installation version %s differs from source %s; restore its source before retrying", name, prior.Version, app.Manifest.Version)
			}
			p.installing[name] = true
		case corevalues.AppStatusActive:
		case corevalues.AppStatusDisabled:
			after = corevalues.AppStatusDisabled
		default:
			return p, fmt.Errorf("App %s has unsupported lifecycle state %q", name, prior.Before)
		}
		if name == "core" && after == corevalues.AppStatusDisabled {
			return p, fmt.Errorf("Core cannot be disabled")
		}
		if after == corevalues.AppStatusActive {
			for _, dep := range app.Manifest.Dependencies {
				if states[dep].Before == corevalues.AppStatusDisabled {
					return p, fmt.Errorf("App %s requires disabled App %s; migration will not reactivate it", name, dep)
				}
			}
		}
		p.Apps = append(p.Apps, AppChange{Name: name, Version: app.Manifest.Version, Before: prior.Before, After: after})
	}
	if _, ok := live.Tables["patch_run"]; ok {
		if _, ok := live.Tables["app"]; !ok {
			return p, fmt.Errorf("partial Core schema: patch_run exists without app")
		}
		// PatchLedger needs only the same read methods, plus Exec for its writer API.
		ledger, ok := q.(db.PatchLedgerQueryer)
		if !ok {
			return p, fmt.Errorf("migration queryer cannot read patch ledger")
		}
		p.runs, err = db.NewPatchLedger(ledger).ListPatchRuns(ctx)
		if err != nil {
			return p, err
		}
	}
	p.loadedPatches, err = patches.Discover(metadata.Apps)
	if err != nil {
		return p, err
	}
	p.PreSync, err = db.BuildLifecyclePatchPlan(p.loadedPatches, metadata.Entities, live, p.runs, db.PatchPhasePreSync, p.installing)
	if err != nil {
		return p, err
	}
	p.Schema, err = db.BuildMetadataSchemaPlan(metadata.Entities, p.PreSync.SchemaAfter)
	if err != nil {
		return p, err
	}
	p.SchemaDeferred = p.PreSync.UnsimulatedSQL
	future, err := db.SchemaAfterSync(metadata.Entities, p.PreSync.SchemaAfter)
	if err != nil {
		return p, err
	}
	p.PostSync, err = db.BuildLifecyclePatchPlan(p.loadedPatches, metadata.Entities, future, p.runs, db.PatchPhasePostSync, p.installing)
	if err != nil {
		return p, err
	}
	// Disabled Apps retain storage and metadata but contribute no runtime access or Fixtures.
	effective := project.Metadata{Entities: metadata.Entities, Pages: metadata.Pages}
	for i, app := range metadata.Apps {
		if p.Apps[i].After == corevalues.AppStatusActive {
			effective.Apps = append(effective.Apps, app)
		}
	}
	var roles []string
	if _, ok := live.Tables["role"]; ok {
		roles, err = access.RoleNames(ctx, q)
		if err != nil {
			return p, err
		}
	}
	p.Access, err = access.Discover(root, effective)
	if err != nil {
		return p, err
	}
	if err := access.ValidateWithPages(&p.Access, metadata.Entities, metadata.Pages, roles); err != nil {
		return p, err
	}
	files, err := fixtures.Discover(effective.Apps)
	if err != nil {
		return p, err
	}
	if err := fixtures.ValidateFiles(files, metadata.Entities); err != nil {
		return p, err
	}
	p.Fixtures = fixtures.Plan{Files: files, Entities: metadata.Entities}
	p.fingerprint, err = fingerprint(root, p)
	return p, err
}

// Apply revalidates the approved plan under the same lock used by schema prune.
// Every migration database effect commits together. Retries start from the last
// committed state; diagnostic Logs and the post-commit snapshot are separate.
func (r Runner) Apply(ctx context.Context, root, databaseURL, version string, approved Plan) (Result, error) {
	pool, err := db.OpenRuntimePool(ctx, databaseURL)
	if err != nil {
		return Result{}, err
	}
	defer pool.Close()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Result{}, err
	}
	defer tx.Rollback(ctx)
	if err := db.LockLifecycle(ctx, tx); err != nil {
		return Result{}, err
	}
	plan, err := r.plan(ctx, tx, root)
	if err != nil {
		return Result{}, err
	}
	if plan.fingerprint != approved.fingerprint {
		return Result{}, fmt.Errorf("migration plan changed; rerun dygo db migrate to review it")
	}
	var names []string
	for _, app := range plan.Apps {
		if plan.installing[app.Name] {
			names = append(names, app.Name)
		}
	}
	ctx = db.WithInstallingApps(ctx, names)
	result := Result{}
	pre := plan.PreSync
	pre.Baseline = nil
	result.PreSync, err = db.ApplyPatchPlanTx(ctx, tx, pre, root, version)
	if err != nil {
		return result, err
	}
	live, err := db.InspectLiveSchema(ctx, tx)
	if err != nil {
		return result, err
	}
	schema, err := db.BuildMetadataSchemaPlan(plan.metadata.Entities, live)
	if err != nil {
		return result, err
	}
	result.Schema, err = db.ApplyMetadataTx(ctx, tx, schema, plan.metadata)
	if err != nil {
		return result, err
	}
	bases := append(plan.PreSync.Baseline, plan.PostSync.Baseline...)
	baselined, err := db.ApplyPatchPlanTx(ctx, tx, db.PatchPlan{Baseline: bases}, root, version)
	if err != nil {
		return result, err
	}
	var postBaselines []db.PatchRun
	for _, run := range baselined.Baselined {
		if run.Phase == db.PatchPhasePostSync {
			postBaselines = append(postBaselines, run)
		} else {
			result.PreSync.Baselined = append(result.PreSync.Baselined, run)
		}
	}
	live, err = db.InspectLiveSchema(ctx, tx)
	if err != nil {
		return result, err
	}
	// Replan post-sync against actual transactional schema; SQL escapes cannot be simulated.
	runs, err := db.NewPatchLedger(tx).ListPatchRuns(ctx)
	if err != nil {
		return result, err
	}
	post, err := db.BuildLifecyclePatchPlan(plan.loadedPatches, plan.metadata.Entities, live, runs, db.PatchPhasePostSync, plan.installing)
	if err != nil {
		return result, err
	}
	result.PostSync, err = db.ApplyPatchPlanTx(ctx, tx, post, root, version)
	if err != nil {
		return result, err
	}
	result.PostSync.Baselined = append(result.PostSync.Baselined, postBaselines...)
	result.Access, err = access.ApplyPlanTx(ctx, tx, plan.Access)
	if err != nil {
		return result, fmt.Errorf("apply access: %w", err)
	}
	result.Fixtures, err = fixtures.NewRunnerWithHooks(r.RecordHooks).ApplyPlanTx(ctx, tx, plan.Fixtures)
	if err != nil {
		return result, fmt.Errorf("apply Fixtures: %w", err)
	}
	for _, app := range plan.Apps {
		// Registry bootstrap defines the metadata that normal Records depend on.
		if _, err := tx.Exec(ctx, `UPDATE "app" SET status=$2, version=$3, updated_at=now() WHERE name=$1`, app.Name, app.After, app.Version); err != nil {
			return result, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return result, err
	}
	if err := r.Migrator.DumpSchema(ctx, root, databaseURL); err != nil {
		return result, fmt.Errorf("migration committed; schema snapshot refresh failed: %w", err)
	}
	return result, nil
}

func fingerprint(root string, plan Plan) ([32]byte, error) {
	h := sha256.New()
	// Source bytes catch changed secrets in Fixture inputs without printing them.
	for _, app := range plan.metadata.Apps {
		err := filepath.WalkDir(app.Dir, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				if entry.Name() == "node_modules" || entry.Name() == ".git" || entry.Name() == "dist" {
					return filepath.SkipDir
				}
				return nil
			}
			switch filepath.Ext(path) {
			case ".yml", ".yaml", ".go":
			default:
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fmt.Fprintln(h, filepath.ToSlash(path))
			h.Write(data)
			return nil
		})
		if err != nil {
			return [32]byte{}, err
		}
	}
	queuesPath := filepath.Join(root, "config", "queues.yml")
	if queues, err := os.ReadFile(queuesPath); err == nil {
		fmt.Fprintln(h, filepath.ToSlash(queuesPath))
		h.Write(queues)
	} else if !os.IsNotExist(err) {
		return [32]byte{}, err
	}
	data, err := json.Marshal(plan)
	if err != nil {
		return [32]byte{}, err
	}
	h.Write(data)
	var sum [32]byte
	copy(sum[:], h.Sum(nil))
	return sum, nil
}
