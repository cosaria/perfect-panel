package access

import (
	accessapi "github.com/perfect-panel/server-v2/internal/domains/access/api"
	accessmodel "github.com/perfect-panel/server-v2/internal/domains/access/model"
)

const (
	RoleAdmin = accessmodel.RoleAdmin
	RoleUser  = accessmodel.RoleUser

	PermissionAdminUsersRead  = accessmodel.PermissionAdminUsersRead
	PermissionAdminUsersWrite = accessmodel.PermissionAdminUsersWrite
)

type RoleSeed = accessmodel.RoleSeed
type PermissionSeed = accessmodel.PermissionSeed
type AdminHTTPHandler = accessapi.AdminHTTPHandler

func RequiredRoles() []RoleSeed {
	return accessmodel.RequiredRoles()
}

func RequiredPermissions() []PermissionSeed {
	return accessmodel.RequiredPermissions()
}

func RequiredRolePermissions() map[string][]string {
	return accessmodel.RequiredRolePermissions()
}

func NewAdminHTTPHandler() *AdminHTTPHandler {
	return accessapi.NewAdminHTTPHandler()
}
