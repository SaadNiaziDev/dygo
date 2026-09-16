package db

import (
	"context"
	"fmt"
	"slices"
)

type installingAppsKey struct{}

// WithInstallingApps grants database lifecycle operations access to newly
// installing Apps before activation. Runtime request handlers must never add
// this scope. It does not change App status.
func WithInstallingApps(ctx context.Context, names []string) context.Context {
	return context.WithValue(ctx, installingAppsKey{}, slices.Clone(names))
}

// AppRuntimePredicate restricts an App table to active Apps, plus explicitly
// scoped lifecycle access. alias must be a trusted SQL table alias, never input.
// Any lifecycle names are appended as a bound parameter to args.
func AppRuntimePredicate(ctx context.Context, alias string, args *[]any) string {
	prefix := ""
	if alias != "" {
		prefix = alias + "."
	}
	active := prefix + "status = 'active'"
	names, _ := ctx.Value(installingAppsKey{}).([]string)
	if len(names) == 0 {
		return active
	}
	*args = append(*args, names)
	return fmt.Sprintf("(%s OR %sname = ANY($%d::text[]))", active, prefix, len(*args))
}
