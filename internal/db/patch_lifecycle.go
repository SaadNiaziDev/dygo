package db

import (
	"context"
	"fmt"
	"time"

	"github.com/hapyco/dygo/internal/entity/catalog"
	"github.com/hapyco/dygo/internal/patches"
	"github.com/jackc/pgx/v5"
)

// BuildLifecyclePatchPlan selects fresh-install work without replaying historical
// migrations against the current schema. All recorded checksums remain immutable.
func BuildLifecyclePatchPlan(loaded []patches.LoadedPatch, entities []catalog.LoadedEntity, live LiveSchema, runs []PatchRun, phase string, installing map[string]bool) (PatchPlan, error) {
	recorded := map[string]PatchRun{}
	for _, run := range runs {
		recorded[patchRunKey(run.AppName, run.PatchID)] = run
	}
	var selected []patches.LoadedPatch
	var baseline []PlannedPatch
	for _, patch := range loaded {
		if patch.Patch.Phase != phase {
			continue
		}
		if run, ok := recorded[patchRunKey(patch.AppName, patch.Patch.ID)]; ok {
			if run.Checksum != patch.Checksum {
				return PatchPlan{}, PatchRunChecksumMismatchError{AppName: patch.AppName, PatchID: patch.Patch.ID, AppliedChecksum: run.Checksum, CurrentChecksum: patch.Checksum}
			}
			selected = append(selected, patch)
			continue
		}
		mode := patch.Patch.RunOn
		if mode == "" {
			mode = patches.RunOnMigrate
		}
		if installing[patch.AppName] && mode == patches.RunOnMigrate {
			baseline = append(baseline, plannedPatchFromLoaded(patch))
			continue
		}
		if !installing[patch.AppName] && mode == patches.RunOnInstall {
			continue
		}
		selected = append(selected, patch)
	}
	plan, err := BuildPatchPlan(selected, entities, live, runs, phase)
	plan.Baseline = baseline
	return plan, err
}

// ApplyPatchPlanTx applies work in its caller's transaction, including the ledger.
// Baselines must be written after schema metadata has established App identities.
func ApplyPatchPlanTx(ctx context.Context, tx pgx.Tx, plan PatchPlan, root, version string) (PatchApplyResult, error) {
	result := PatchApplyResult{Phase: plan.Phase}
	for _, group := range []struct {
		patches []PlannedPatch
		outcome string
	}{{plan.Pending, PatchOutcomeApplied}, {plan.Baseline, PatchOutcomeBaselined}} {
		for _, patch := range group.patches {
			if group.outcome == PatchOutcomeApplied {
				for _, operation := range patch.Operations {
					if err := executePatchOperation(ctx, tx, patch, operation); err != nil {
						return result, fmt.Errorf("apply patch %s/%s: %w", patch.AppName, patch.PatchID, err)
					}
				}
			}
			path, err := patchLedgerPath(root, patch)
			if err != nil {
				return result, err
			}
			run := PatchRun{AppName: patch.AppName, PatchID: patch.PatchID, Path: path, Phase: patch.Phase, Checksum: patch.Checksum, AppliedAt: time.Now().UTC(), DygoVersion: version, Outcome: group.outcome}
			if err := NewPatchLedger(tx).RecordPatchRun(ctx, run); err != nil {
				return result, err
			}
			if group.outcome == PatchOutcomeBaselined {
				result.Baselined = append(result.Baselined, run)
			} else {
				result.Applied = append(result.Applied, run)
			}
		}
	}
	return result, nil
}

// SchemaAfterSync provides the expected columns for read-only post-sync Patch
// planning. Apply always reinspects the transaction after executing schema changes.
func SchemaAfterSync(entities []catalog.LoadedEntity, live LiveSchema) (LiveSchema, error) {
	desired, err := buildDesiredSchema(entities)
	if err != nil {
		return LiveSchema{}, err
	}
	result := cloneLiveSchema(live)
	for _, table := range desired.Tables {
		current, ok := result.Tables[table.Name]
		if !ok {
			current = liveTable{Name: table.Name, Columns: map[string]liveColumn{}}
		}
		for _, col := range append(table.SystemColumns, table.Columns...) {
			current.Columns[col.Name] = liveColumn{Name: col.Name, Type: col.Type, Nullable: !col.Required, HasDefault: col.HasSafeDefault, DefaultSQL: col.DefaultSQL}
		}
		result.Tables[table.Name] = current
	}
	return result, nil
}
