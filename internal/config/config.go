// internal/config/config.go
package config

import (
	"fmt"
	"os"
)

var (
	JWTSecret  string
	AuthHeader string
)

func Init() {
	JWTSecret = os.Getenv("JWT_SECRET")
	AuthHeader = os.Getenv("AUTH_HEADER")

	// Valores por defecto
	if JWTSecret == "" {
		JWTSecret = "default_jwt_secret"
	}
	fmt.Println("JWTSecret cargado:", JWTSecret) // Para debug
	if AuthHeader == "" {
		AuthHeader = "dev-secret-header"
	}
}
