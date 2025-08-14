// internal/utils/jwt.go
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("default_jwt_secret") // Usa variable de entorno en producción

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT() (string, error) {
	claims := Claims{
		UserID: 123,
		Role:   "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "your-app.com",
			Audience:  jwt.ClaimStrings{"your-client-id"},
			Subject:   "user@example.com",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
