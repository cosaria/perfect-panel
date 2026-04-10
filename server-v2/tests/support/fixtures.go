package support

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// WriteJSONFixture 将对象写入临时 JSON fixture 文件并返回路径。
func WriteJSONFixture(t *testing.T, dir string, fileName string, value any) string {
	t.Helper()

	content, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("序列化 fixture 失败: %v", err)
	}

	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("写入 fixture 失败: %v", err)
	}

	return path
}
