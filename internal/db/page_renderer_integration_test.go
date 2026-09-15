package db

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/hapyco/dygo/internal/patches"
)

func TestPageRendererUpgradePreservesRowsAndRollsBack(t *testing.T) {
	pool, _ := auditDatabase(t)
	ctx := context.Background()
	_, err := pool.Exec(ctx, `CREATE TABLE page (renderer text NOT NULL CONSTRAINT page_renderer_check CHECK (renderer IN ('entity-index'))); INSERT INTO page VALUES ('entity-index')`)
	if err != nil {
		t.Fatal(err)
	}
	patch, err := patches.LoadFile(filepath.Join("..", "..", "apps", "core", "patches", "0002_enable_app_page_renderers.yml"))
	if err != nil {
		t.Fatal(err)
	}
	statement := patch.Operations[0].Fields["statement"].Value
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, statement); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO page VALUES ('vue')`); err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO page VALUES ('vue')`); err == nil {
		t.Fatal("rollback must restore the old constraint")
	}
	for range 2 {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := pool.Exec(ctx, `INSERT INTO page VALUES ('vue')`); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM page WHERE renderer = 'entity-index'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("existing Pages count=%d, err=%v", count, err)
	}
}
