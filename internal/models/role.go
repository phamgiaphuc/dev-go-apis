package models

type Role struct {
	ID          int     `db:"id" json:"id" binding:"required"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}

type RoleIDs []int

type RolePermissions struct {
	Role
	Permissions   Permissions   `db:"permissions" json:"-"`
	PermissionIDs PermissionIDs `json:"permission_ids"`
}

type RolePermissionsList []RolePermissions

type GetRolePermissionsByIdRequest struct {
	ID int `uri:"id" binding:"required"`
}

type CreateRoleRequest struct {
	Name        string  `db:"name" json:"name" binding:"required"`
	Description *string `db:"description" json:"description"`
}

type UpdateRolePermissionsRequest struct {
	RolePermissions `json:",inline"`
}

type DeleteRolePermissionsRequest struct {
	RoleIDs `json:"role_ids" binding:"required"`
}
