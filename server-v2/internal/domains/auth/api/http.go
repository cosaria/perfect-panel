package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	authmodel "github.com/perfect-panel/server-v2/internal/domains/auth/model"
	authusecase "github.com/perfect-panel/server-v2/internal/domains/auth/usecase"
)

type HTTPHandler struct {
	service *authusecase.Service
}

func NewHTTPHandler(service *authusecase.Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var input authmodel.SignInInput
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := h.service.SignIn(r.Context(), input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	WriteEnvelope(w, http.StatusCreated, map[string]any{
		"id":          result.Session.ID,
		"userId":      result.Session.UserID,
		"accessToken": result.AccessToken,
		"expiresAt":   result.Session.ExpiresAt.UTC().Format(time.RFC3339),
	})
}

func (h *HTTPHandler) CreateVerificationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := h.service.IssueVerificationToken(r.Context(), authmodel.IssueVerificationInput{
		Email:   input.Email,
		Purpose: authmodel.VerificationPurposeEmailVerification,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	WriteEnvelope(w, http.StatusCreated, map[string]any{
		"id":          result.Token.ID,
		"destination": result.Token.Destination,
		"expiresAt":   result.Token.ExpiresAt.UTC().Format(time.RFC3339),
	})
}

func (h *HTTPHandler) CreatePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var input authmodel.RequestPasswordResetInput
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := h.service.RequestPasswordReset(r.Context(), input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	WriteEnvelope(w, http.StatusCreated, map[string]any{
		"id":         result.Token.ID,
		"status":     "accepted",
		"acceptedAt": h.service.Clock().Now().UTC().Format(time.RFC3339),
	})
}

func (h *HTTPHandler) CreatePasswordReset(w http.ResponseWriter, r *http.Request) {
	var input authmodel.ResetPasswordInput
	if !decodeJSON(w, r, &input) {
		return
	}
	result, err := h.service.ResetPassword(r.Context(), input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	WriteEnvelope(w, http.StatusOK, map[string]any{
		"id":          result.TokenID,
		"status":      result.Status,
		"completedAt": result.CompletedAt.UTC().Format(time.RFC3339),
	})
}

func (h *HTTPHandler) ListMySessions(w http.ResponseWriter, r *http.Request) {
	principal, ok := authmodel.PrincipalFromContext(r.Context())
	if !ok {
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "缺少会话上下文", "unauthorized")
		return
	}
	sessions, err := h.service.ListUserSessions(r.Context(), principal.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	data := make([]map[string]any, 0, len(sessions))
	for _, session := range sessions {
		data = append(data, map[string]any{
			"id":        session.ID,
			"userId":    session.UserID,
			"createdAt": session.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	WriteEnvelopeWithMeta(w, http.StatusOK, data, map[string]any{
		"page":       1,
		"pageSize":   len(data),
		"total":      len(data),
		"totalPages": 1,
		"hasNext":    false,
		"hasPrev":    false,
	})
}

func (h *HTTPHandler) DeleteMySession(w http.ResponseWriter, r *http.Request) {
	principal, ok := authmodel.PrincipalFromContext(r.Context())
	if !ok {
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "缺少会话上下文", "unauthorized")
		return
	}
	sessionID := strings.TrimSpace(r.PathValue("sessionId"))
	if sessionID == "" {
		WriteProblem(w, http.StatusBadRequest, "Bad Request", "缺少 sessionId 路径参数", "bad_request")
		return
	}
	if err := h.service.SignOut(r.Context(), authmodel.SignOutInput{
		UserID:    principal.UserID,
		SessionID: sessionID,
	}); err != nil {
		writeServiceError(w, err)
		return
	}
	WriteEnvelope(w, http.StatusOK, map[string]any{
		"revoked":   true,
		"revokedAt": h.service.Clock().Now().UTC().Format(time.RFC3339),
	})
}

type envelope struct {
	Data any            `json:"data"`
	Meta map[string]any `json:"meta"`
}

type problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
	Code   string `json:"code,omitempty"`
}

type validationProblem struct {
	Type   string            `json:"type"`
	Title  string            `json:"title"`
	Status int               `json:"status"`
	Detail string            `json:"detail"`
	Code   string            `json:"code,omitempty"`
	Errors []validationError `json:"errors"`
}

type validationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteEnvelope(w http.ResponseWriter, status int, data any) {
	WriteEnvelopeWithMeta(w, status, data, map[string]any{})
}

func WriteEnvelopeWithMeta(w http.ResponseWriter, status int, data any, meta map[string]any) {
	writeJSON(w, status, envelope{
		Data: data,
		Meta: meta,
	})
}

func WriteProblem(w http.ResponseWriter, status int, title string, detail string, code string) {
	writeJSON(w, status, problem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
		Code:   code,
	})
}

func WriteValidationProblem(w http.ResponseWriter, status int, title string, detail string, code string, errors []validationError) {
	writeJSON(w, status, validationProblem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
		Code:   code,
		Errors: errors,
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authmodel.ErrInvalidCredentials):
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", err.Error(), "invalid_credentials")
	case errors.Is(err, authmodel.ErrUnauthorized), errors.Is(err, authmodel.ErrSessionNotFound):
		WriteProblem(w, http.StatusUnauthorized, "Unauthorized", err.Error(), "unauthorized")
	case errors.Is(err, authmodel.ErrForbidden):
		WriteProblem(w, http.StatusForbidden, "Forbidden", err.Error(), "forbidden")
	case errors.Is(err, authmodel.ErrVerificationTokenConsumed):
		WriteProblem(w, http.StatusConflict, "Conflict", err.Error(), "token_consumed")
	case errors.Is(err, authmodel.ErrVerificationTokenExpired), errors.Is(err, authmodel.ErrVerificationTokenPurposeMismatch):
		WriteValidationProblem(w, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "validation_failed", []validationError{
			{Field: "token", Code: "invalid", Message: err.Error()},
		})
	case errors.Is(err, authmodel.ErrIdentityNotFound), errors.Is(err, authmodel.ErrUserNotFound), errors.Is(err, authmodel.ErrVerificationTokenNotFound):
		WriteValidationProblem(w, http.StatusUnprocessableEntity, "Unprocessable Entity", err.Error(), "validation_failed", []validationError{
			{Field: "email", Code: "not_found", Message: err.Error()},
		})
	default:
		WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", err.Error(), "internal_error")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		WriteProblem(w, http.StatusBadRequest, "Bad Request", "请求体不是合法 JSON", "bad_request")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
