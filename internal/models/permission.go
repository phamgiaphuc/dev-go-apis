package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	DashboardGroup     = "dashboard:group"
	UserDashboardGroup = "dashboard:user"
)

var (
	ReadUserDashboard   = ReadPermission(UserDashboardGroup)
	CreateUserDashboard = CreatePermission(UserDashboardGroup)
	EditUserDashboard   = EditPermission(UserDashboardGroup)
	DeleteUserDashboard = DeletePermission(UserDashboardGroup)
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

type Permission struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	GroupID     int    `db:"group_id" json:"-"`
}

type PermissionGroup struct {
	ID          int    `db:"id" json:"id" binding:"required"`
	Name        string `db:"name" json:"name" binding:"required"`
	Description string `db:"description" json:"description"`
}

type Permissions []Permission

type PermissionList struct {
	PermissionGroup
	Permissions Permissions `db:"permissions" json:"permissions"`
}

func (p *Permissions) Scan(src interface{}) error {
	if src == nil {
		*p = []Permission{}
		return nil
	}

	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("Permissions: type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, p)
}

func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}
