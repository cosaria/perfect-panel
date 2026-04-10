package health

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data any            `json:"data"`
	Meta map[string]any `json:"meta"`
}

// NewHandler 返回健康检查 HTTP 处理器。
func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(envelope{
			Data: map[string]string{
				"status": "ok",
			},
			Meta: map[string]any{},
		})
	})
}
