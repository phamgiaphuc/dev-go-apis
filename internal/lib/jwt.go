package lib

import (
	"dev-go-apis/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ACCESS_TOKEN_TTL_DURATION  = ParseTimeDuration(ACCESS_TOKEN_TTL, 7, time.Hour*24)
	REFRESH_TOKEN_TTL_DURATION = ParseTimeDuration(REFRESH_TOKEN_TTL, 7, time.Hour*24)
)

func SignAccessToken(payload *models.UserWithClaims) (string, error) {
	now := time.Now()
	claims := models.UserWithClaims{
		User:      payload.User,
		SessionID: payload.SessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ACCESS_TOKEN_TTL_DURATION)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(ACCESS_TOKEN_SECRET))
}

func SignRefreshToken(payload *models.UserWithClaims) (string, error) {
	now := time.Now()
	claims := models.UserWithClaims{
		User:      payload.User,
		SessionID: payload.SessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(REFRESH_TOKEN_TTL_DURATION)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(REFRESH_TOKEN_SECRET))
}
