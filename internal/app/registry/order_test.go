package registry

import (
	"strings"
	"testing"

	"github.com/hapyco/dygo/internal/app/manifest"
)

func TestDependencyOrder(t *testing.T) {
	apps := []manifest.LoadedApp{
		loadedApp("sales", "crm"),
		loadedApp("core"),
		loadedApp("crm", "core"),
	}
	ordered, err := DependencyOrder(apps)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, app := range ordered {
		names = append(names, app.Manifest.Name)
	}
	if got := strings.Join(names, ","); got != "core,crm,sales" {
		t.Fatalf("DependencyOrder() = %s", got)
	}
}

func TestDependencyOrderRejectsInvalidGraphs(t *testing.T) {
	for _, test := range []struct {
		name string
		apps []manifest.LoadedApp
		want string
	}{
		{"missing", []manifest.LoadedApp{loadedApp("sales", "crm")}, `depends on unknown app "crm"`},
		{"cycle", []manifest.LoadedApp{loadedApp("crm", "sales"), loadedApp("sales", "crm")}, "dependency cycle"},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := DependencyOrder(test.apps)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("DependencyOrder() error = %v, want %q", err, test.want)
			}
		})
	}
}

func loadedApp(name string, dependencies ...string) manifest.LoadedApp {
	return manifest.LoadedApp{
		ManifestPath: name + "/app.yml",
		Manifest:     manifest.Manifest{Name: name, Dependencies: dependencies},
	}
}
