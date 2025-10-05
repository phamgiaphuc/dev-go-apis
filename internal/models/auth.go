package models

import (
	"github.com/golang-jwt/jwt/v5"
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

type JwtVerification struct {
	Verification
	jwt.RegisteredClaims
}

type JwtTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User       User   `json:"user"`
	Token      string `json:"token,omitempty"`
	IsVerified bool   `json:"is_verified"`
}

type RegisterResponse struct {
	User       User   `json:"user"`
	Token      string `json:"token,omitempty"`
	IsVerified bool   `json:"is_verified"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
