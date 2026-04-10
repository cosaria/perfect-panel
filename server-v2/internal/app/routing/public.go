package routing

import (
	"net/http"

	authdomain "github.com/perfect-panel/server-v2/internal/domains/auth"
	authapi "github.com/perfect-panel/server-v2/internal/domains/auth/api"
)

func RegisterPublic(mux *http.ServeMux, handler *authdomain.HTTPHandler) {
	if mux == nil || handler == nil {
		return
	}
	authapi.RegisterPublicRoutes(mux, handler)
}
