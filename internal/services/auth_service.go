// internal/services/auth_service.go
package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/repositories"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(user *domain.User) error {
	return s.userRepo.Create(user)
}

// ... Login, etc.
