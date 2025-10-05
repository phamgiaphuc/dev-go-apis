package models

import (
	"time"

	"github.com/google/uuid"
)

type Identifier string

const (
	EmailVerifiedIdentifier Identifier = "email-verified"
	PasswordResetIdentifier Identifier = "password-reset"
)

type Verification struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	Identifier Identifier `db:"identifier" json:"identifier"`
	Value      string     `db:"value" json:"value"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	ExpiredAt  time.Time  `db:"expired_at" json:"expired_at"`
}

type CreateVerification struct {
	Verification Verification
	Token        string
	Code         string
}
