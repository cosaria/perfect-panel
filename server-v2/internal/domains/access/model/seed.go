package model

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	PermissionAdminUsersRead  = "admin.users.read"
	PermissionAdminUsersWrite = "admin.users.write"
)

type RoleSeed struct {
	Code string
	Name string
}

type PermissionSeed struct {
	Code string
	Name string
}

func RequiredRoles() []RoleSeed {
	return []RoleSeed{
		{Code: RoleAdmin, Name: "管理员"},
		{Code: RoleUser, Name: "普通用户"},
	}
}

func RequiredPermissions() []PermissionSeed {
	return []PermissionSeed{
		{Code: PermissionAdminUsersRead, Name: "查看后台用户"},
		{Code: PermissionAdminUsersWrite, Name: "写入后台用户"},
	}
}

func RequiredRolePermissions() map[string][]string {
	return map[string][]string{
		RoleAdmin: {
			PermissionAdminUsersRead,
			PermissionAdminUsersWrite,
		},
	}
}
