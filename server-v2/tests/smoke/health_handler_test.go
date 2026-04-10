package smoke_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/perfect-panel/server-v2/internal/platform/http/health"
)

func TestHealthHandlerReturnsEnvelope(t *testing.T) {
	handler := health.NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("状态码应为 200，实际: %d", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应 JSON 解析失败: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("响应缺少 data 对象，body=%s", rec.Body.String())
	}
	if data["status"] != "ok" {
		t.Fatalf("data.status 应为 ok，实际: %v", data["status"])
	}

	meta, ok := body["meta"].(map[string]any)
	if !ok {
		t.Fatalf("响应缺少 meta 对象，body=%s", rec.Body.String())
	}
	if len(meta) != 0 {
		t.Fatalf("meta 应为空对象，实际: %v", meta)
	}
}
