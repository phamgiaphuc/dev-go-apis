package models

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderCredential Provider = "credential"
	ProviderGoogle     Provider = "google"
	ProviderGithub     Provider = "github"
)

type Account struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"-"`
	AccountID  string    `db:"account_id" json:"account_id"`
	ProviderId Provider  `db:"provider_id" json:"provider"`
	Password   string    `db:"password" json:"-"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type GoogleAccountInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}
