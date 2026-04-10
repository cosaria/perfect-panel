package routing

import (
	"net/http"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
	authapi "github.com/perfect-panel/server-v2/internal/domains/auth/api"
)

func RegisterUser(mux *http.ServeMux, handler *authdomain.HTTPHandler, requireSession func(http.Handler) http.Handler) {
	if mux == nil || handler == nil || requireSession == nil {
		return
	}
	authapi.RegisterUserRoutes(mux, handler, requireSession)
}
