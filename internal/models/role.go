package models

type Role struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

type RoleList struct {
	Role
	Permissions Permissions `db:"permissions" json:"permissions"`
}
