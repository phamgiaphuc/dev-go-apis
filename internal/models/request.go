package models

type CreatePermissionGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionGroupRequest struct {
	PermissionGroup
	PermissionIDs []int `json:"permission_ids"`
}
