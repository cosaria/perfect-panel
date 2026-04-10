package api

import "net/http"

func RegisterAdminRoutes(
	mux *http.ServeMux,
	handler *AdminHTTPHandler,
	requireSession func(http.Handler) http.Handler,
	requirePermissions func(...string) func(http.Handler) http.Handler,
	readPermission string,
	writePermission string,
) {
	if mux == nil || handler == nil || requireSession == nil || requirePermissions == nil {
		return
	}

	mux.Handle(
		"GET /api/v1/admin/users",
		requireSession(requirePermissions(readPermission)(http.HandlerFunc(handler.ListUsers))),
	)
	mux.Handle(
		"POST /api/v1/admin/users",
		requireSession(requirePermissions(writePermission)(http.HandlerFunc(handler.CreateUser))),
	)
}
