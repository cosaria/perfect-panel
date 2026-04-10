package contract_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestContractPipelinePasses(t *testing.T) {
	moduleRoot := getModuleRoot(t)

	cmd := exec.Command("make", "contract")
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make contract 失败: %v\n输出:\n%s", err, string(out))
	}

	bundlePath := filepath.Join(moduleRoot, "openapi", "dist", "openapi.json")
	if info, err := os.Stat(bundlePath); err != nil {
		t.Fatalf("未生成 bundle 文件 %s: %v", bundlePath, err)
	} else if info.IsDir() {
		t.Fatalf("bundle 路径是目录而不是文件: %s", bundlePath)
	}

	generatedPath := filepath.Join(moduleRoot, "openapi", "generated", "ts")
	info, err := os.Stat(generatedPath)
	if err != nil {
		t.Fatalf("未生成 openapi-ts 输出目录 %s: %v", generatedPath, err)
	}
	if !info.IsDir() {
		t.Fatalf("openapi-ts 输出路径不是目录: %s", generatedPath)
	}

	entries, err := os.ReadDir(generatedPath)
	if err != nil {
		t.Fatalf("读取 openapi-ts 输出目录失败: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("openapi-ts 输出目录为空: %s", generatedPath)
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
