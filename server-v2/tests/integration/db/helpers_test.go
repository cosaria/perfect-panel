package db_test

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func newIsolatedPostgres(t *testing.T) (dbDSN string, cleanup func()) {
	t.Helper()

	baseDSN := strings.TrimSpace(os.Getenv("TEST_POSTGRES_DSN"))
	if baseDSN == "" {
		baseDSN = strings.TrimSpace(os.Getenv("PPANEL_DB_DSN"))
	}
	if baseDSN == "" {
		t.Skip("未设置 TEST_POSTGRES_DSN 或 PPANEL_DB_DSN，跳过 Postgres 集成测试")
	}

	baseDB, err := ppdb.Open(baseDSN)
	if err != nil {
		t.Fatalf("连接测试数据库失败: %v", err)
	}

	schema := fmt.Sprintf("it_%08x", rand.Uint32())
	if _, err := baseDB.Exec(`CREATE SCHEMA IF NOT EXISTS "` + schema + `"`); err != nil {
		_ = baseDB.Close()
		t.Fatalf("创建测试 schema 失败: %v", err)
	}

	dsnWithSchema, err := applySearchPath(baseDSN, schema)
	if err != nil {
		_ = baseDB.Close()
		t.Fatalf("构造 search_path DSN 失败: %v", err)
	}

	cleanup = func() {
		_, _ = baseDB.Exec(`DROP SCHEMA IF EXISTS "` + schema + `" CASCADE`)
		_ = baseDB.Close()
	}
	return dsnWithSchema, cleanup
}

func applySearchPath(dsn string, schema string) (string, error) {
	if strings.Contains(dsn, "://") {
		parsed, err := url.Parse(dsn)
		if err != nil {
			return "", err
		}
		query := parsed.Query()
		query.Set("search_path", schema)
		parsed.RawQuery = query.Encode()
		return parsed.String(), nil
	}
	return dsn + " search_path=" + schema, nil
}

func createSchemaWithTargetRevisionButContractDrifted(t *testing.T, dsn string) {
	t.Helper()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	statements := []string{
		`CREATE TABLE schema_revisions (
			id BIGSERIAL PRIMARY KEY,
			version TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`INSERT INTO schema_revisions(version) VALUES ('` + ppdb.TargetSchemaVersion + `')`,
		`CREATE TABLE users (
			id BIGINT PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE roles (
			id UUID PRIMARY KEY,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE system_settings (
			id UUID PRIMARY KEY,
			scope TEXT NOT NULL,
			key TEXT NOT NULL,
			value_json JSONB NOT NULL DEFAULT '{}'::jsonb,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE UNIQUE INDEX idx_system_settings_scope_key ON system_settings(scope, key)`,
		`CREATE TABLE outbox_events (
			id UUID PRIMARY KEY,
			topic TEXT NOT NULL,
			status TEXT NOT NULL,
			aggregate_type TEXT NOT NULL,
			aggregate_id UUID NOT NULL,
			payload JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("准备漂移 schema 失败: %v", err)
		}
	}
}

func createSchemaWithTargetRevisionMissingIDDefaults(t *testing.T, dsn string) {
	t.Helper()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	statements := []string{
		`CREATE TABLE schema_revisions (
			id BIGSERIAL PRIMARY KEY,
			version TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`INSERT INTO schema_revisions(version) VALUES ('` + ppdb.TargetSchemaVersion + `')`,
		`CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE roles (
			id UUID PRIMARY KEY,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE system_settings (
			id UUID PRIMARY KEY,
			scope TEXT NOT NULL,
			key TEXT NOT NULL,
			value_json JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE UNIQUE INDEX idx_system_settings_scope_key ON system_settings(scope, key)`,
		`CREATE TABLE outbox_events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			topic TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			aggregate_type TEXT NOT NULL,
			aggregate_id UUID NOT NULL,
			payload JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("准备缺省默认值 schema 失败: %v", err)
		}
	}
}

func createBaselineOnlySchema(t *testing.T, dsn string) {
	t.Helper()

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	scriptPath := filepath.Join(moduleRoot(t), "internal", "platform", "db", "migrations", "0001_baseline.sql")
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("读取 baseline migration 失败: %v", err)
	}
	if _, err := db.Exec(string(script)); err != nil {
		t.Fatalf("执行 baseline migration 失败: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO schema_revisions(version) VALUES ('0001_baseline')`); err != nil {
		t.Fatalf("写入 baseline revision 失败: %v", err)
	}
}
