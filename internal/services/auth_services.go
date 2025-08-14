// internal/services/auth_service.go
package services

import (
	"errors"
	"os"
)

type AuthService struct {
	// Puedes agregar dependencias como userRepo si necesitas
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Authenticate(username, password string) (bool, error) {
	envUser := os.Getenv("ADMIN_USERNAME")
	envPass := os.Getenv("ADMIN_PASSWORD")

	if envUser == "" || envPass == "" {
		return false, errors.New("credenciales de administrador no configuradas")
	}

	return username == envUser && password == envPass, nil
}
