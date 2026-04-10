package middleware

import (
	"net/http"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
)

func RequirePermissions(required ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := authdomain.PrincipalFromContext(r.Context())
			if !ok {
				authdomain.WriteProblem(w, http.StatusUnauthorized, "Unauthorized", "缺少会话上下文", "unauthorized")
				return
			}
			for _, permission := range required {
				if !hasPermission(principal.Permissions, permission) {
					authdomain.WriteProblem(w, http.StatusForbidden, "Forbidden", "权限不足", "forbidden")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func hasPermission(existing []string, required string) bool {
	for _, permission := range existing {
		if permission == required {
			return true
		}
	}
	return false
}
