package db_test

import (
	"context"
	"testing"

	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func TestMigrateFailsOnDriftedSchema(t *testing.T) {
	dsn, cleanup := newIsolatedPostgres(t)
	defer cleanup()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE users (id BIGINT PRIMARY KEY)`); err != nil {
		t.Fatalf("预置漂移 users 表失败: %v", err)
	}

	err = ppdb.Migrate(context.Background(), db)
	if err == nil {
		t.Fatal("存在漂移表时，migrate 应失败")
	}

	var revisionTableExists bool
	if err := db.QueryRow(`SELECT to_regclass(current_schema() || '.schema_revisions') IS NOT NULL`).Scan(&revisionTableExists); err != nil {
		t.Fatalf("检查 schema_revisions 是否存在失败: %v", err)
	}
	if !revisionTableExists {
		return
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_revisions WHERE version = $1`, ppdb.TargetSchemaVersion).Scan(&count); err != nil {
		t.Fatalf("查询目标版本 revision 失败: %v", err)
	}
	if count != 0 {
		t.Fatalf("漂移失败场景不应写入目标版本 revision，实际 count=%d", count)
	}
}
