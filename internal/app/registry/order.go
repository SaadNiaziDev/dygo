package registry

import (
	"fmt"
	"github.com/hapyco/dygo/internal/app/manifest"
	"sort"
	"strings"
)

// DependencyOrder returns Apps in deterministic dependency order, rejecting missing dependencies and cycles.
func DependencyOrder(apps []manifest.LoadedApp) ([]manifest.LoadedApp, error) {
	byName := map[string]manifest.LoadedApp{}
	for _, app := range apps {
		name := app.Manifest.Name
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("app name is required")
		}
		if previous, ok := byName[name]; ok {
			return nil, fmt.Errorf("duplicate app %q in %s and %s", name, previous.ManifestPath, app.ManifestPath)
		}
		byName[name] = app
	}

	indegree := map[string]int{}
	dependents := map[string][]string{}
	for _, app := range apps {
		name := app.Manifest.Name
		indegree[name] = 0
	}
	for _, app := range apps {
		name := app.Manifest.Name
		seenDependencies := map[string]struct{}{}
		for _, dependency := range app.Manifest.Dependencies {
			if _, ok := byName[dependency]; !ok {
				return nil, fmt.Errorf("app %q depends on unknown app %q", name, dependency)
			}
			if _, ok := seenDependencies[dependency]; ok {
				continue
			}
			seenDependencies[dependency] = struct{}{}
			indegree[name]++
			dependents[dependency] = append(dependents[dependency], name)
		}
	}
	for dependency := range dependents {
		sort.Strings(dependents[dependency])
	}

	var ready []string
	for name, degree := range indegree {
		if degree == 0 {
			ready = append(ready, name)
		}
	}
	sort.Strings(ready)

	ordered := make([]manifest.LoadedApp, 0, len(apps))
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		ordered = append(ordered, byName[name])
		for _, dependent := range dependents[name] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sort.Strings(ready)
	}
	if len(ordered) != len(apps) {
		var cycle []string
		for name, degree := range indegree {
			if degree > 0 {
				cycle = append(cycle, name)
			}
		}
		sort.Strings(cycle)
		return nil, fmt.Errorf("app dependency cycle among %s", strings.Join(cycle, ", "))
	}
	return ordered, nil
}
