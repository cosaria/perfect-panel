package smoke_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCliBootstrapFilesExist(t *testing.T) {
	moduleRoot := getModuleRoot(t)

	requiredFiles := []string{
		"cmd/server/main.go",
		"cmd/server/root.go",
		"internal/app/runtime/modes.go",
	}

	for _, rel := range requiredFiles {
		abs := filepath.Join(moduleRoot, rel)
		if _, err := os.Stat(abs); err != nil {
			t.Fatalf("文件不存在: %s, 错误: %v", rel, err)
		}
	}
}

func TestCliHelpBoots(t *testing.T) {
	moduleRoot := getModuleRoot(t)

	cmd := exec.Command("go", "run", "./cmd/server", "--help")
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run --help 失败: %v, 输出: %s", err, string(out))
	}

	output := string(out)
	if !strings.Contains(output, "serve-api") {
		t.Fatalf("帮助信息未包含 serve-api 子命令，输出: %s", output)
	}
}

func getModuleRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法获取当前测试文件路径")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
