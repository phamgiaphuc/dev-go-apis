package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

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
