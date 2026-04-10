package db_test

import (
	"testing"

	"github.com/perfect-panel/server-v2/internal/app/bootstrap"
	"github.com/perfect-panel/server-v2/internal/app/runtime"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func TestServeFailsWhenSchemaVersionMismatches(t *testing.T) {
	dsn, cleanup := newIsolatedPostgres(t)
	defer cleanup()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	if err := ppdb.Migrate(t.Context(), db); err != nil {
		t.Fatalf("执行 migrate 失败: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM schema_revisions`); err != nil {
		t.Fatalf("清空 schema_revisions 失败: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO schema_revisions(version) VALUES ('0000_prebaseline')`); err != nil {
		t.Fatalf("插入错误版本失败: %v", err)
	}

	t.Setenv("PPANEL_DB_DSN", dsn)

	_, err = bootstrap.BuildForMode(runtime.ModeServeAPI, bootstrap.Options{})
	if err == nil {
		t.Fatal("schema version 不匹配时，serve 模式构建应失败")
	}
}
