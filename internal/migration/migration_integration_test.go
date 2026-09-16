package migration

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/hapyco/dygo/internal/db"
	"github.com/hapyco/dygo/internal/projectgen"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresGeneratedProjectMigrationLifecycle(t *testing.T) {
	ctx := context.Background()
	databaseURL := migrationDatabase(t)
	root := generatedProject(t, databaseURL)
	runner := Runner{Migrator: db.Migrator{}}

	plan, err := runner.Plan(ctx, root, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Apps) != 3 {
		t.Fatalf("Apps = %+v, want generated Core, Studio, and Business App", plan.Apps)
	}
	if len(plan.PreSync.Baseline) == 0 {
		t.Fatal("fresh installation did not baseline historical patches")
	}

	queuesPath := filepath.Join(root, "config", "queues.yml")
	queues, err := os.ReadFile(queuesPath)
	if err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(string(queues), "concurrency: 4", "concurrency: 5", 1)
	if err := os.WriteFile(queuesPath, []byte(changed), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := runner.Apply(ctx, root, databaseURL, "test", plan); err == nil || !strings.Contains(err.Error(), "migration plan changed") {
		t.Fatalf("Apply() after queue change error = %v", err)
	}
	if err := os.WriteFile(queuesPath, queues, 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := runner.Apply(ctx, root, databaseURL, "test", plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.PreSync.Baselined) != len(plan.PreSync.Baseline) || len(result.PostSync.Baselined) != len(plan.PostSync.Baseline) {
		t.Fatalf("baselined pre/post = %d/%d, want %d/%d", len(result.PreSync.Baselined), len(result.PostSync.Baselined), len(plan.PreSync.Baseline), len(plan.PostSync.Baseline))
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var active, apps, baselined, roles, countries int
	if err := pool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE status='active'), count(*) FROM "app"`).Scan(&active, &apps); err != nil {
		t.Fatal(err)
	}
	if active != apps || apps != 3 {
		t.Fatalf("active Apps = %d/%d", active, apps)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM "patch_run" WHERE outcome='baselined'`).Scan(&baselined); err != nil {
		t.Fatal(err)
	}
	if baselined != len(plan.PreSync.Baseline)+len(plan.PostSync.Baseline) {
		t.Fatalf("baselined patch rows = %d", baselined)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM "role"`).Scan(&roles); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM "country"`).Scan(&countries); err != nil {
		t.Fatal(err)
	}
	if roles == 0 || countries == 0 {
		t.Fatalf("migration did not apply access and Fixtures: roles=%d countries=%d", roles, countries)
	}
	if _, err := pool.Exec(ctx, `UPDATE "app" SET status='disabled' WHERE name='studio'`); err != nil {
		t.Fatal(err)
	}
	disabledPlan, err := runner.Plan(ctx, root, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runner.Apply(ctx, root, databaseURL, "test", disabledPlan); err != nil {
		t.Fatal(err)
	}
	pages, err := db.NewMetadataReader(pool).ListPages(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range pages {
		if page.App.Name == "studio" {
			t.Fatal("disabled Studio page remained runtime-visible")
		}
	}
}

func TestPostgresMigrationRollsBackLateFixtureFailureAndRetries(t *testing.T) {
	ctx := context.Background()
	databaseURL := migrationDatabase(t)
	root := generatedProject(t, databaseURL)
	hooks := db.NewRecordHookRegistry()
	if err := hooks.RegisterGlobal(db.RecordBeforeCreate, "reject-fixture", func(context.Context, db.RecordHookContext) error {
		return errors.New("forced fixture failure")
	}); err != nil {
		t.Fatal(err)
	}
	failing := Runner{Migrator: db.Migrator{}, RecordHooks: hooks}
	plan, err := failing.Plan(ctx, root, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := failing.Apply(ctx, root, databaseURL, "test", plan); err == nil || !strings.Contains(err.Error(), "record hook failed") {
		t.Fatalf("Apply() error = %v", err)
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	var appTable *string
	if err := pool.QueryRow(ctx, `SELECT to_regclass('public.app')::text`).Scan(&appTable); err != nil {
		pool.Close()
		t.Fatal(err)
	}
	pool.Close()
	if appTable != nil {
		t.Fatalf("failed migration left schema behind: %s", *appTable)
	}

	if _, err := (Runner{Migrator: db.Migrator{}}).Apply(ctx, root, databaseURL, "test", plan); err != nil {
		t.Fatalf("retry after rollback: %v", err)
	}
}

func generatedProject(t *testing.T, databaseURL string) string {
	t.Helper()
	repositoryRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	result, err := projectgen.Generate(context.Background(), projectgen.Options{
		Name:          "Acme",
		ModulePath:    "example.com/acme",
		WorkingDir:    t.TempDir(),
		FrameworkRoot: repositoryRoot,
		DatabaseURL:   databaseURL,
		SkipTidy:      true,
		StudioAssets: fstest.MapFS{
			"index.html":    {Data: []byte(`<div id="app"></div>`)},
			"assets/app.js": {Data: []byte(`console.log("studio")`)},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return result.Path
}

func migrationDatabase(t *testing.T) string {
	t.Helper()
	baseURL := os.Getenv("DYGO_TEST_DATABASE_URL")
	if baseURL == "" {
		t.Skip("set DYGO_TEST_DATABASE_URL to run PostgreSQL regressions")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, baseURL)
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("dygo_migration_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" {
		admin.Close()
		t.Fatalf("DYGO_TEST_DATABASE_URL must be a URL: %v", err)
	}
	parsed.Path = "/" + name
	databaseURL := parsed.String()
	t.Cleanup(func() {
		if _, err := admin.Exec(ctx, "DROP DATABASE "+name+" WITH (FORCE)"); err != nil {
			t.Error(err)
		}
		admin.Close()
	})
	return databaseURL
}
