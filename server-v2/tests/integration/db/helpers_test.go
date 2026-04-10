package db_test

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"os"
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
