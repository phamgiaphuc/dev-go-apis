package lib

import (
	"time"

	"dev-go-apis/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ACCESS_TOKEN_TTL_DURATION  = ParseTimeDuration(ACCESS_TOKEN_TTL, 15, time.Minute)
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

func ParseToken(token string, secret string) (*models.UserWithClaims, error) {
	tok, err := jwt.ParseWithClaims(token, &models.UserWithClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*models.UserWithClaims)
	if !ok || !tok.Valid {
		return nil, err
	}
	return claims, nil
}
