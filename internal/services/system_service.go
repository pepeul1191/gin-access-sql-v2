package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/forms"
	"accessv2/internal/repositories"
	"errors"
	"time"
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

func (s *SystemService) GetPaginatedSystems(page, perPage int, nameQuery, descQuery string) ([]domain.System, int64, error) {
	// Validación básica
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Delegar al repositorio
	return s.repo.GetPaginated(page, perPage, nameQuery, descQuery)
}

func (s *SystemService) CreateSystem(input *forms.SystemCreateInput) (*domain.System, error) {

	var (
		ErrSystemNameRequired = errors.New("el nombre del sistema es requerido")
	)
	// Validación de datos
	if input.Name == "" {
		return nil, ErrSystemNameRequired
	}

	// Crear objeto del dominio
	system := &domain.System{
		Name:        input.Name,
		Description: input.Description,
		Repository:  input.Repository,
		Created:     input.Created,
		Updated:     input.Updated,
	}

	// Establecer fechas por defecto si no vienen
	if system.Created.IsZero() {
		system.Created = time.Now()
	}
	if system.Updated.IsZero() {
		system.Updated = system.Created
	}

	// Guardar en la base de datos
	if err := s.repo.Create(system); err != nil {
		return nil, err
	}

	return system, nil
}
