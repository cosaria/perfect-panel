package db_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
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

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法定位测试文件路径")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
