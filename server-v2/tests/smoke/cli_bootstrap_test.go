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
		"cmd/server/serve_api.go",
		"cmd/server/serve_worker.go",
		"cmd/server/serve_scheduler.go",
		"cmd/server/migrate.go",
		"cmd/server/seed_required.go",
		"cmd/server/seed_demo.go",
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
	for _, mode := range allModes() {
		if !strings.Contains(output, mode) {
			t.Fatalf("帮助信息未包含 %s 子命令，输出: %s", mode, output)
		}
	}
}

func TestCliSubCommandsHelpBoots(t *testing.T) {
	moduleRoot := getModuleRoot(t)

	for _, mode := range allModes() {
		cmd := exec.Command("go", "run", "./cmd/server", mode, "--help")
		cmd.Dir = moduleRoot
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s --help 失败: %v, 输出: %s", mode, err, string(out))
		}
	}
}

func TestCliRejectsUnexpectedArgs(t *testing.T) {
	moduleRoot := getModuleRoot(t)

	testCases := [][]string{
		{"migrate", "--bad-flag"},
		{"serve-api", "unexpected-arg"},
	}

	for _, args := range testCases {
		cmd := exec.Command("go", append([]string{"run", "./cmd/server"}, args...)...)
		cmd.Dir = moduleRoot
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("命令应失败但成功: %v, 输出: %s", args, string(out))
		}
		if !strings.Contains(string(out), "不接受额外参数") {
			t.Fatalf("错误信息不符合预期: %v, 输出: %s", args, string(out))
		}
	}
}

func allModes() []string {
	return []string{
		"serve-api",
		"serve-worker",
		"serve-scheduler",
		"migrate",
		"seed-required",
		"seed-demo",
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
