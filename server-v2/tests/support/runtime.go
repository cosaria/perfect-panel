package support

import (
	"path/filepath"
	"runtime"
	"testing"
)

// RuntimeEnv 描述测试运行时的最小环境信息。
type RuntimeEnv struct {
	ModuleRoot string
	TempDir    string
}

// ModuleRoot 返回 server-v2 模块根目录。
func ModuleRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法获取测试支架路径")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

// NewRuntimeEnv 构建测试运行时环境。
func NewRuntimeEnv(t *testing.T) RuntimeEnv {
	t.Helper()

	return RuntimeEnv{
		ModuleRoot: ModuleRoot(t),
		TempDir:    t.TempDir(),
	}
}
