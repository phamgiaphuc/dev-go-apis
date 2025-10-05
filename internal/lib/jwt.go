package lib

import (
	"errors"
	"reflect"
	"time"

	"dev-go-apis/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ACCESS_TOKEN_TTL_DURATION   = ParseTimeDuration(ACCESS_TOKEN_TTL, 15, time.Minute)
	REFRESH_TOKEN_TTL_DURATION  = ParseTimeDuration(REFRESH_TOKEN_TTL, 7, time.Hour*24)
	VERIFIED_EMAIL_TTL_DURATION = time.Minute * 5
)

func SignVerificationToken(payload *models.UserVerification, duration time.Duration) (*models.JwtUserVerificationToken, error) {
	now := time.Now()
	expiredAt := now.Add(duration)
	claims := models.JwtUserVerification{
		UserVerification: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}

	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(VERIFICATION_TOKEN_SECRET))
	if err != nil {
		return nil, err
	}

	return &models.JwtUserVerificationToken{
		Token:     tok,
		ExpiredAt: expiredAt,
	}, nil
}

func SignAccessToken(payload *models.UserWithClaims) (string, error) {
	now := time.Now()
	claims := models.JwtUserPayload{
		UserWithClaims: *payload,
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
	claims := models.JwtUserPayload{
		UserWithClaims: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(REFRESH_TOKEN_TTL_DURATION)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(REFRESH_TOKEN_SECRET))
}

func ParseToken[T jwt.Claims](tokenString string, secret string) (*T, error) {
	if tokenString == "" {
		return nil, errors.New("empty token string")
	}

	claimsPtr := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface().(jwt.Claims)

	token, err := jwt.ParseWithClaims(tokenString, claimsPtr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims type")
	}

	return &claims, nil
}
