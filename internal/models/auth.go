package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type VerificationIdentifier string

const (
	EmailVerifiedIdentifier VerificationIdentifier = "email-verified"
	PasswordResetIdentifier VerificationIdentifier = "password-reset"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type JwtUserPayload struct {
	UserWithClaims
	jwt.RegisteredClaims
}

type JwtUserVerification struct {
	UserVerification
	jwt.RegisteredClaims
}

type JwtUserVerificationToken struct {
	Token     string
	ExpiredAt time.Time
}

type JwtTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User       User      `json:"user"`
	Token      string    `json:"token"`
	IsVerified bool      `json:"is_verified"`
	ExpiredAt  time.Time `json:"expired_at,omitempty"`
}

type RegisterResponse struct {
	User       `json:"user"`
	Token      string    `json:"token"`
	IsVerified bool      `json:"is_verified"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
