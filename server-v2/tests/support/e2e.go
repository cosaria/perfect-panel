package support

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewE2EServer 启动最小 HTTP 测试服务器。
func NewE2EServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}
