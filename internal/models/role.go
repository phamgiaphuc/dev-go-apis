package models

import "github.com/lib/pq"

type Role struct {
	ID          int     `db:"id" json:"id" binding:"required"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}

type RoleIDs []int

type JSONInt64Array pq.Int64Array
type JSONStringArray pq.StringArray

type RoleWithPermissions struct {
	Role          `json:",inline"`
	PermissionIDs pq.Int64Array  `db:"permission_ids" json:"permission_ids" swaggertype:"array,integer"`
	Permissions   pq.StringArray `db:"permissions" json:"-"`
}

type RoleWithPermissions2 struct {
	Role          `json:",inline"`
	PermissionIDs []int64  `json:"permission_ids" swaggertype:"array,integer"`
	Permissions   []string `json:"permissions"`
}

type RoleList []RoleWithPermissions

type GetRoleByIdRequest struct {
	ID int `uri:"id" binding:"required"`
}

type CreateRoleRequest struct {
	Name        string  `db:"name" json:"name" binding:"required"`
	Description *string `db:"description" json:"description"`
}

type UpdateRoleRequest struct {
	RoleWithPermissions `json:",inline"`
}

type DeleteRolesRequest struct {
	RoleIDs `json:"role_ids" binding:"required"`
}
