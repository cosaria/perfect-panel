package db_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
)

func TestSeedCommandsFailWhenSchemaContractDrifted(t *testing.T) {
	for _, command := range []string{"seed-required", "seed-demo"} {
		t.Run(command, func(t *testing.T) {
			dsn, cleanup := newIsolatedPostgres(t)
			defer cleanup()

			createSchemaWithTargetRevisionButContractDrifted(t, dsn)
			cmd := exec.Command("go", "run", "./cmd/server", command)
			cmd.Dir = moduleRoot(t)
			cmd.Env = append(os.Environ(), "PPANEL_DB_DSN="+dsn)

			if output, err := cmd.CombinedOutput(); err == nil {
				t.Fatalf("schema 契约漂移时，%s 应失败", command)
			} else if len(output) == 0 {
				t.Fatalf("%s 失败时应输出错误信息", command)
			}
		})
	}
}

func TestSeedCommandsFailWhenIDDefaultMissing(t *testing.T) {
	for _, command := range []string{"seed-required", "seed-demo"} {
		t.Run(command, func(t *testing.T) {
			dsn, cleanup := newIsolatedPostgres(t)
			defer cleanup()

			createSchemaWithTargetRevisionMissingIDDefaults(t, dsn)
			cmd := exec.Command("go", "run", "./cmd/server", command)
			cmd.Dir = moduleRoot(t)
			cmd.Env = append(os.Environ(), "PPANEL_DB_DSN="+dsn)

			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("id 默认值缺失时，%s 应失败", command)
			}
			if !strings.Contains(strings.ToLower(string(output)), "schema 契约") {
				t.Fatalf("%s 应在 gate 阶段失败，实际输出: %s", command, string(output))
			}
		})
	}
}

func TestSeedRequiredDoesNotApplyAuthAccessMigrationImplicitly(t *testing.T) {
	dsn, cleanup := newIsolatedPostgres(t)
	defer cleanup()

	createBaselineOnlySchema(t, dsn)

	cmd := exec.Command("go", "run", "./cmd/server", "seed-required")
	cmd.Dir = moduleRoot(t)
	cmd.Env = append(os.Environ(), "PPANEL_DB_DSN="+dsn)

	if output, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("baseline-only schema 上，seed-required 不应隐式执行 0002 migration，实际输出: %s", string(output))
	}

	db, err := ppdb.Open(dsn)
	if err != nil {
		t.Fatalf("连接隔离数据库失败: %v", err)
	}
	defer db.Close()

	var permissionsTableExists bool
	if err := db.QueryRow(`SELECT to_regclass(current_schema() || '.permissions') IS NOT NULL`).Scan(&permissionsTableExists); err != nil {
		t.Fatalf("检查 permissions 表是否存在失败: %v", err)
	}
	if permissionsTableExists {
		t.Fatal("seed-required 不应通过隐式迁移创建 permissions 表")
	}

	var revisionCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_revisions`).Scan(&revisionCount); err != nil {
		t.Fatalf("统计 schema_revisions 失败: %v", err)
	}
	if revisionCount != 1 {
		t.Fatalf("seed-required 不应改写 revision 主链，want=1 got=%s", strconv.Itoa(revisionCount))
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位测试文件路径")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}
