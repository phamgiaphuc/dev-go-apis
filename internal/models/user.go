package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"email_verified"`
	Image         string    `db:"image" json:"image"`
	IsBanned      bool      `db:"is_banned" json:"is_banned"`
	RoleID        int       `db:"role_id" json:"role_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`

	Role string `json:"role"`
}

type UserWithOAuth struct {
	User    User
	Account Account
}

type UserWithAccount struct {
	User    User    `json:"user" db:"user"`
	Account Account `json:"account" db:"account"`
}

type UserWithClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	jwt.RegisteredClaims
}

type UserWithAccounts struct {
	User
	Accounts []Account `json:"accounts"`
}

type GetMeResponse struct {
	User      User      `json:"user"`
	SessionID uuid.UUID `json:"session_id"`
}
