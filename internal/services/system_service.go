package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/repositories"
)

type SystemService struct {
	repo *repositories.SystemRepository
}

func NewSystemService(repo *repositories.SystemRepository) *SystemService {
	return &SystemService{repo: repo}
}

func (s *SystemService) GetAllSystems() ([]domain.System, error) {
	return s.repo.GetAll()
}
