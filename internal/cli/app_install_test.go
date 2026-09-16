package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppInstallPreparesDependencyClosureAndRunner(t *testing.T) {
	root := t.TempDir()
	writeCLIProjectRoot(t, root)
	writeCLIGoModule(t, root, "example.com/acme")
	writeCLIApp(t, filepath.Join(root, ".dygo", "apps", "core"), "core")
	writeCLIAppWithBody(t, filepath.Join(root, "apps", "crm"), `
name: crm
label: CRM
version: 0.1.0
dependencies: [core]
`)
	writeCLIAppWithBody(t, filepath.Join(root, "apps", "sales"), `
name: sales
label: Sales
version: 0.1.0
dependencies: [crm]
`)
	t.Chdir(root)

	var stdout bytes.Buffer
	if err := Run(context.Background(), []string{"app", "install", "sales", "--yes"}, strings.NewReader(""), &stdout, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"prepare order: core -> crm -> sales", "Next: build the project runner, then run dygo db migrate."} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout = %q, want %q", stdout.String(), want)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "cmd", "dygo", "main.go")); err != nil {
		t.Fatalf("runner was not written: %v", err)
	}
}

func TestAppInstallDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	writeCLIProjectRoot(t, root)
	writeCLIGoModule(t, root, "example.com/acme")
	writeCLIApp(t, filepath.Join(root, "apps", "sales"), "sales")
	t.Chdir(root)

	var stdout bytes.Buffer
	if err := Run(context.Background(), []string{"app", "install", "sales", "--dry-run"}, strings.NewReader(""), &stdout, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "dry-run: no files or database changes") {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if _, err := os.Stat(filepath.Join(root, "cmd", "dygo", "main.go")); !os.IsNotExist(err) {
		t.Fatalf("dry-run runner stat = %v, want missing", err)
	}
}

func TestAppInstallRejectsUnknownApp(t *testing.T) {
	root := t.TempDir()
	writeCLIProjectRoot(t, root)
	writeCLIGoModule(t, root, "example.com/acme")
	t.Chdir(root)

	err := Run(context.Background(), []string{"app", "install", "missing", "--yes"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), `app "missing" was not found`) {
		t.Fatalf("error = %v", err)
	}
}
