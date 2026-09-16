package db

import (
	"context"
	"reflect"
	"testing"
)

func TestAppRuntimePredicateScopesOnlyInstallingApps(t *testing.T) {
	args := []any{"existing"}
	if got := AppRuntimePredicate(context.Background(), "a", &args); got != "a.status = 'active'" {
		t.Fatalf("active predicate = %q", got)
	}
	ctx := WithInstallingApps(context.Background(), []string{"crm", "sales"})
	if got := AppRuntimePredicate(ctx, "a", &args); got != "(a.status = 'active' OR a.name = ANY($2::text[]))" {
		t.Fatalf("install predicate = %q", got)
	}
	if !reflect.DeepEqual(args[1], []string{"crm", "sales"}) {
		t.Fatalf("install names = %#v", args[1])
	}
}
