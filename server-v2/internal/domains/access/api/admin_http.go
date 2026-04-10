package api

import (
	"encoding/json"
	"net/http"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
)

type AdminHTTPHandler struct{}

func NewAdminHTTPHandler() *AdminHTTPHandler {
	return &AdminHTTPHandler{}
}

func (h *AdminHTTPHandler) ListUsers(w http.ResponseWriter, _ *http.Request) {
	authdomain.WriteEnvelopeWithMeta(w, http.StatusOK, []map[string]any{}, map[string]any{
		"page":       1,
		"pageSize":   0,
		"total":      0,
		"totalPages": 0,
		"hasNext":    false,
		"hasPrev":    false,
	})
}

func (h *AdminHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email  string `json:"email"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		authdomain.WriteProblem(w, http.StatusBadRequest, "Bad Request", "请求体不是合法 JSON", "bad_request")
		return
	}
	authdomain.WriteEnvelope(w, http.StatusCreated, map[string]any{
		"id":     "00000000-0000-4000-8000-000000000001",
		"email":  input.Email,
		"status": input.Status,
	})
}
