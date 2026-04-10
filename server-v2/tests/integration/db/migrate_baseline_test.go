package db_test

import (
	"context"
	"fmt"
	"testing"

	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func TestMigrateAppliesBaselineTables(t *testing.T) {
	dsn, cleanup := newIsolatedPostgres(t)
	defer cleanup()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	if err := ppdb.Migrate(context.Background(), db); err != nil {
		t.Fatalf("执行 migrate 失败: %v", err)
	}

	tables := []string{"users", "roles", "system_settings", "outbox_events", "schema_revisions"}
	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = current_schema() AND table_name = $1
		)`
		if err := db.QueryRow(query, table).Scan(&exists); err != nil {
			t.Fatalf("查询表 %s 失败: %v", table, err)
		}
		if !exists {
			t.Fatalf("预期表 %s 存在，但不存在", table)
		}
	}

	version, err := ppdb.CurrentSchemaVersion(context.Background(), db)
	if err != nil {
		t.Fatalf("读取 schema version 失败: %v", err)
	}
	if version != ppdb.TargetSchemaVersion {
		t.Fatalf("schema version 不匹配，want=%s got=%s", ppdb.TargetSchemaVersion, version)
	}

	var revisions int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_revisions WHERE version = $1`, ppdb.TargetSchemaVersion).Scan(&revisions); err != nil {
		t.Fatalf("统计 schema_revisions 失败: %v", err)
	}
	if revisions == 0 {
		t.Fatal("migrate 后未写入目标 schema version")
	}

	t.Log(fmt.Sprintf("baseline migration 已就绪，version=%s", version))
}
