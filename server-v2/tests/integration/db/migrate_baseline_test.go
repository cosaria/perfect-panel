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

	tables := []string{
		"users",
		"roles",
		"system_settings",
		"outbox_events",
		"schema_revisions",
		"user_identities",
		"user_sessions",
		"verification_tokens",
		"permissions",
		"role_permissions",
		"user_roles",
		"auth_events",
	}
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

	var currentRevisionCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_revisions WHERE version = $1`, ppdb.TargetSchemaVersion).Scan(&currentRevisionCount); err != nil {
		t.Fatalf("统计目标 revision 失败: %v", err)
	}
	if currentRevisionCount == 0 {
		t.Fatal("migrate 后未写入目标 schema version")
	}

	var baselineRevisionCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_revisions WHERE version = '0001_baseline'`).Scan(&baselineRevisionCount); err != nil {
		t.Fatalf("统计 baseline revision 失败: %v", err)
	}
	if baselineRevisionCount == 0 {
		t.Fatal("migrate 后应保留 baseline revision 记录")
	}

	t.Log(fmt.Sprintf("baseline migration 已就绪，version=%s", version))
}

func TestMigrateUpgradesBaselineRevisionToAuthAccess(t *testing.T) {
	dsn, cleanup := newIsolatedPostgres(t)
	defer cleanup()

	createBaselineOnlySchema(t, dsn)

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	if err := ppdb.Migrate(context.Background(), db); err != nil {
		t.Fatalf("baseline -> auth/access migrate 失败: %v", err)
	}

	version, err := ppdb.CurrentSchemaVersion(context.Background(), db)
	if err != nil {
		t.Fatalf("读取 schema version 失败: %v", err)
	}
	if version != ppdb.TargetSchemaVersion {
		t.Fatalf("schema version 不匹配，want=%s got=%s", ppdb.TargetSchemaVersion, version)
	}

	var permissionsTableExists bool
	if err := db.QueryRow(`SELECT to_regclass(current_schema() || '.permissions') IS NOT NULL`).Scan(&permissionsTableExists); err != nil {
		t.Fatalf("检查 permissions 表是否存在失败: %v", err)
	}
	if !permissionsTableExists {
		t.Fatal("升级到 auth/access revision 后，permissions 表必须存在")
	}

	var revisionCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_revisions WHERE version IN ('0001_baseline', '0002_auth_access')`).Scan(&revisionCount); err != nil {
		t.Fatalf("统计 revision 链失败: %v", err)
	}
	if revisionCount != 2 {
		t.Fatalf("revision 主链应包含 0001 与 0002，got=%d", revisionCount)
	}
}
