package lib

import (
	"fmt"
)

const (
	AdminDashboard     = "dashboard:admin"
	UserAdminDashboard = "dashboard:user"
	RoleAdminDashboard = "dashboard:role"
)

var (
	ReadUserAdminDashboard   = ReadPermission(UserAdminDashboard)
	CreateUserAdminDashboard = CreatePermission(UserAdminDashboard)
	EditUserAdminDashboard   = EditPermission(UserAdminDashboard)
	DeleteUserAdminDashboard = DeletePermission(UserAdminDashboard)

	ReadRoleAdminDashboard   = ReadPermission(RoleAdminDashboard)
	CreateRoleAdminDashboard = CreatePermission(RoleAdminDashboard)
	EditRoleAdminDashboard   = EditPermission(RoleAdminDashboard)
	DeleteRoleAdminDashboard = ReadPermission(RoleAdminDashboard)

	TempAdminDashboard = ReadPermission("dashboard:temp")
)

func ReadPermission(group string) string {
	return fmt.Sprintf("%s:read", group)
}

func CreatePermission(group string) string {
	return fmt.Sprintf("%s:create", group)
}

func EditPermission(group string) string {
	return fmt.Sprintf("%s:edit", group)
}

func DeletePermission(group string) string {
	return fmt.Sprintf("%s:delete", group)
}

var (
	// User APIs
	GetUserByIdPermissions = []string{
		ReadUserAdminDashboard,
		CreateUserAdminDashboard,
		EditUserAdminDashboard,
		DeleteUserAdminDashboard,
	}
)
