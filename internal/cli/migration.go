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
	"github.com/spf13/cobra"
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
					return fmt.Errorf("App %s Entity %s Hooks are not compiled; run dygo app install %s and build the project runner", entity.AppName, entity.Entity.Name, entity.AppName)
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
				return fmt.Errorf("App %s Job %s is not compiled; run dygo app install %s and build the project runner", job.AppName, job.Job.Name, job.AppName)
			}
		}
		return nil
	}
}

// Keep the former spelling discoverable without allowing a second database writer.
func retiredApplyCommand(resource string) *cobra.Command {
	cmd := &cobra.Command{Use: "apply", Short: "Use dygo db migrate to apply database changes", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		return fmt.Errorf("dygo %s apply has moved into dygo db migrate; run dygo db migrate --env %s", resource, cmd.Flag("env").Value.String())
	}}
	cmd.Flags().String("env", "development", "Environment for the replacement database command")
	cmd.Flags().Bool("yes", false, "Use --yes with dygo db migrate")
	cmd.Flags().Bool("dry-run", false, "Use --dry-run with dygo db migrate")
	return cmd
}
