package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hapyco/dygo/internal/app/manifest"
	"github.com/hapyco/dygo/internal/app/registry"
	"github.com/hapyco/dygo/internal/project"
	"github.com/hapyco/dygo/internal/runnergen"
	"github.com/spf13/cobra"
)

func newAppInstallCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	var dryRun, yes bool
	cmd := &cobra.Command{
		Use:   "install <app>",
		Short: "Validate a local App and prepare its project runner",
		Long: "Validate locally available Apps and dependencies, then update the generated project runner. " +
			"App source must already exist under apps/ or .dygo/apps/. " +
			"Build the project runner and run dygo db migrate to apply database changes.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := workingRootPath()
			if err != nil {
				return err
			}
			if err := runnergen.RequireGeneratedProjectRoot(root); err != nil {
				return fmt.Errorf("app install requires a generated dygo project: %w", err)
			}
			// Share the existing app validate command's validation, without invoking a subprocess.
			apps, err := project.LoadApps(root)
			if err != nil {
				return err
			}
			order, err := appInstallOrder(apps, args[0])
			if err != nil {
				return err
			}
			// Wiring always covers the whole project, independently of database state.
			update, err := runnergen.Render(root, runnergen.RenderOptions{})
			if err != nil {
				return fmt.Errorf("plan app installation: %w", err)
			}
			before, err := os.ReadFile(update.RunnerFile)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("read project runner: %w", err)
			}
			missing := os.IsNotExist(err)
			status, err := runnergen.GeneratedFileStatus(update.RunnerFile, update.Source, true)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(stdout, "app install plan\nproject: %s\n%d apps are valid\nprepare order: %s\nrunner: %s (%s; all local Apps)\n", root, len(apps), strings.Join(order, " -> "), relToHooksRoot(root, update.RunnerFile), status); err != nil {
				return fmt.Errorf("write app install plan: %w", err)
			}
			if dryRun {
				_, err := fmt.Fprintln(stdout, "dry-run: no files or database changes")
				return err
			}
			if status != "unchanged" {
				if !yes {
					confirmed, err := streamConfirmer(stdin, stderr)(cmd.Context(), "Apply local App installation?")
					if err != nil {
						return err
					}
					if !confirmed {
						_, err := fmt.Fprintln(stdout, "App installation cancelled.")
						return err
					}
				}
				// Reject edits made while the user reviewed the plan; write only its rendered source.
				current, err := runnergen.Render(root, runnergen.RenderOptions{})
				if err != nil {
					return err
				}
				existing, err := os.ReadFile(update.RunnerFile)
				if err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("read project runner: %w", err)
				}
				if !bytes.Equal(current.Source, update.Source) || !bytes.Equal(before, existing) || missing != os.IsNotExist(err) {
					return fmt.Errorf("project runner plan changed; rerun dygo app install %s", args[0])
				}
				if err := os.MkdirAll(filepath.Dir(update.RunnerFile), 0o755); err != nil {
					return fmt.Errorf("create project runner directory: %w", err)
				}
				if _, err := runnergen.WriteFileIfChanged(update.RunnerFile, update.Source); err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(stdout, "Apps prepared: %s\nNext: build the project runner, then run dygo db migrate.\n", strings.Join(order, ", "))
			return err
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show local installation changes without writing")
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip the local installation confirmation")
	return cmd
}

func appInstallOrder(apps []manifest.LoadedApp, target string) ([]string, error) {
	ordered, err := registry.DependencyOrder(apps)
	if err != nil {
		return nil, err
	}
	wanted := map[string]bool{target: true}
	found := false
	for index := len(ordered) - 1; index >= 0; index-- {
		app := ordered[index]
		if app.Manifest.Name == target {
			found = true
		}
		if wanted[app.Manifest.Name] {
			for _, dependency := range app.Manifest.Dependencies {
				wanted[dependency] = true
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("app %q was not found; add its source under apps/ or .dygo/apps/ before installation", target)
	}
	var names []string
	for _, app := range ordered {
		if wanted[app.Manifest.Name] {
			names = append(names, app.Manifest.Name)
		}
	}
	return names, nil
}
