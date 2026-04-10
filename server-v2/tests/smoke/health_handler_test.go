package smoke_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/perfect-panel/server-v2/internal/platform/http/health"
	"github.com/perfect-panel/server-v2/tests/support"
)

func TestHealthHandlerReturnsEnvelope(t *testing.T) {
	srv := support.NewE2EServer(t, health.NewHandler())
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("调用健康检查失败: %v", err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("响应 JSON 解析失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("状态码应为 200，实际: %d", resp.StatusCode)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("响应缺少 data 对象")
	}
	if data["status"] != "ok" {
		t.Fatalf("data.status 应为 ok，实际: %v", data["status"])
	}

	meta, ok := body["meta"].(map[string]any)
	if !ok {
		t.Fatalf("响应缺少 meta 对象")
	}
	if len(meta) != 0 {
		t.Fatalf("meta 应为空对象，实际: %v", meta)
	}
}
