package middleware

import (
	"net/http"
	"strings"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
)

func RequireSession(resolver authdomain.SessionResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerTokenFromHeader(r.Header.Get("Authorization"))
			if !ok {
				authdomain.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "缺少 Bearer Token", "unauthorized")
				return
			}
			principal, err := resolver.ResolveSession(r.Context(), token)
			if err != nil {
				authdomain.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", err.Error(), "unauthorized")
				return
			}
			next.ServeHTTP(w, r.WithContext(authdomain.WithPrincipal(r.Context(), principal)))
		})
	}
}

func bearerTokenFromHeader(header string) (string, bool) {
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(header)), "bearer ") {
		return "", false
	}
	token := strings.TrimSpace(header[len("Bearer "):])
	return token, token != ""
}
