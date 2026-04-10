package support

import (
	"path/filepath"
	"runtime"
	"testing"
)

// ModuleRoot 返回 server-v2 模块根目录。
func ModuleRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("无法获取测试支架路径")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
