package routing

import (
	"net/http"

	"github.com/perfect-panel/server-v2/internal/domains/access"
	accessapi "github.com/perfect-panel/server-v2/internal/domains/access/api"
)

func RegisterAdmin(
	mux *http.ServeMux,
	handler *access.AdminHTTPHandler,
	requireSession func(http.Handler) http.Handler,
	requirePermissions func(...string) func(http.Handler) http.Handler,
) {
	if mux == nil || handler == nil || requireSession == nil || requirePermissions == nil {
		return
	}
	accessapi.RegisterAdminRoutes(mux, handler, requireSession, requirePermissions, access.PermissionAdminUsersRead, access.PermissionAdminUsersWrite)
}
