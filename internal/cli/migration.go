package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hapyco/dygo/internal/db"
	jobruntime "github.com/hapyco/dygo/internal/jobs/runtime"
	"github.com/hapyco/dygo/internal/migration"
	"github.com/hapyco/dygo/internal/project"
)

type projectMigrator struct {
	db.Migrator
	lifecycle migration.Runner
}

func (m projectMigrator) MigrationPlan(ctx context.Context, root, url string) (migration.Plan, error) {
	return m.lifecycle.Plan(ctx, root, url)
}
func (m projectMigrator) Migrate(ctx context.Context, root, url string, plan migration.Plan) (migration.Result, error) {
	return m.lifecycle.Apply(ctx, root, url, currentVersion(), plan)
}

func migrationCodeValidator(hooks *db.RecordHookRegistry, jobs *jobruntime.Registry) func(string, project.RuntimeMetadata) error {
	return func(root string, metadata project.RuntimeMetadata) error {
		for _, entity := range metadata.Entities {
			path := filepath.Join(filepath.Dir(entity.Path), "hooks.go")
			if _, err := os.Stat(path); err == nil {
				if !hooks.HasEntity(entity.AppName, entity.Entity.Name) {
					return fmt.Errorf("App %s Entity %s Hooks are not compiled; run dygo hook sync and build the project runner", entity.AppName, entity.Entity.Name)
				}
			} else if !os.IsNotExist(err) {
				return err
			}
		}
		for _, job := range metadata.Jobs {
			// The worker binds this Core handler to the environment mailer at startup.
			if job.AppName == "core" && job.Job.Name == "send-notification-email" {
				continue
			}
			if !jobs.HasJob(job.AppName, job.Job.Name) {
				return fmt.Errorf("App %s Job %s is not compiled; run dygo hook sync and build the project runner", job.AppName, job.Job.Name)
			}
		}
		return nil
	}
}
