package api

import "net/http"

func RegisterPublicRoutes(mux *http.ServeMux, handler *HTTPHandler) {
	if mux == nil || handler == nil {
		return
	}

	mux.Handle("POST /api/v1/public/sessions", http.HandlerFunc(handler.CreateSession))
	mux.Handle("POST /api/v1/public/verification-tokens", http.HandlerFunc(handler.CreateVerificationToken))
	mux.Handle("POST /api/v1/public/password-reset-requests", http.HandlerFunc(handler.CreatePasswordResetRequest))
	mux.Handle("POST /api/v1/public/password-resets", http.HandlerFunc(handler.CreatePasswordReset))
}

func RegisterUserRoutes(mux *http.ServeMux, handler *HTTPHandler, requireSession func(http.Handler) http.Handler) {
	if mux == nil || handler == nil || requireSession == nil {
		return
	}

	mux.Handle("GET /api/v1/user/me/sessions", requireSession(http.HandlerFunc(handler.ListMySessions)))
	mux.Handle("DELETE /api/v1/user/me/sessions/{sessionId}", requireSession(http.HandlerFunc(handler.DeleteMySession)))
}
