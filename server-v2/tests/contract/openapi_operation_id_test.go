package contract_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOperationIDsAreExplicit(t *testing.T) {
	moduleRoot := getModuleRoot(t)
	bundlePath := filepath.Join(moduleRoot, "openapi", "dist", "openapi.json")

	runContractPipeline(t, moduleRoot)

	raw, err := os.ReadFile(bundlePath)
	if err != nil {
		t.Fatalf("读取 bundle 失败: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("解析 bundle 失败: %v", err)
	}

	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		t.Fatalf("bundle 中缺少 paths 节点")
	}

	seen := make(map[string]string)
	var operationCount int
	for path, rawItem := range paths {
		item, ok := rawItem.(map[string]any)
		if !ok {
			t.Fatalf("路径 %s 的定义格式不正确", path)
		}

		for method, rawOperation := range item {
			if !isHTTPMethod(method) {
				continue
			}

			operationCount++
			operation, ok := rawOperation.(map[string]any)
			if !ok {
				t.Fatalf("路径 %s 的 %s 操作格式不正确", path, method)
			}

			operationId, ok := operation["operationId"].(string)
			if !ok || strings.TrimSpace(operationId) == "" {
				t.Fatalf("路径 %s 的 %s 操作缺少显式 operationId", path, method)
			}
			if prev, exists := seen[operationId]; exists {
				t.Fatalf("operationId 重复: %s，已被 %s 和 %s 同时使用", operationId, prev, fmt.Sprintf("%s %s", path, method))
			}
			seen[operationId] = fmt.Sprintf("%s %s", path, method)
		}
	}

	if operationCount == 0 {
		t.Fatalf("bundle 中没有发现任何 operation")
	}
}

func isHTTPMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "post", "put", "patch", "delete", "options", "head", "trace":
		return true
	default:
		return false
	}
}
