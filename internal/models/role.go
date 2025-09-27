package models

type Role struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}

type RoleList struct {
	Role
	Permissions Permissions `db:"permissions" json:"permissions"`
}

type CreateRoleRequest struct {
	Name        string  `db:"name" json:"name" binding:"required"`
	Description *string `db:"description" json:"description"`
}

type GetRoleListResponse struct {
	Roles []RoleList `json:"roles"`
}

type CreateRoleResponse struct {
	Role *Role `json:"role"`
}
