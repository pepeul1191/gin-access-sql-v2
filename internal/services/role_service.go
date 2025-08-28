package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/forms"
	"accessv2/internal/repositories"
	"errors"
	"time"
)

type RoleService struct {
	repo *repositories.RoleRepository
}

func NewRoleService(repo *repositories.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAllBySystemID(systemID int) ([]domain.Role, error) {
	// Delegar al repositorio
	return s.repo.GetRolesBySystemID(systemID)
}

func (s *RoleService) CreateRole(input *forms.RoleCreateInput, systemID int) (*domain.Role, error) {

	var (
		ErrUserNameRequired = errors.New("El nombre del rol es requerido")
	)
	// Validación de datos
	if input.Name == "" {
		return nil, ErrUserNameRequired
	}

	// Validación única
	err := s.repo.CheckRoleExistsInSystem(input.Name, systemID)
	if err != nil {
		return nil, errors.New("Nombre de sol ya en uso en el sistema")
	}
	// Crear objeto del dominio
	role := &domain.Role{
		Name: input.Name,
	}

	// Establecer fechas por defecto si no vienen
	if role.Created.IsZero() {
		role.Created = time.Now()
	}
	if role.Updated.IsZero() {
		role.Updated = role.Created
	}

	// Guardar en la base de datos
	if err := s.repo.Create(role); err != nil {
		return nil, err
	}

	return role, nil
}

// FetchSystem usando el repository (versión con pointer)
func (s *RoleService) FetchRole(id uint64, role *domain.Role) error {
	if role == nil {
		return errors.New("Role pointer cannot be nil")
	}

	// Usar el método del repository que acepta pointer
	tempRole, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Copiar los valores al User proporcionado
	*role = tempRole
	return nil
}

// UpdateRole usando el repository
func (s *RoleService) UpdateRole(role *domain.Role) error {
	if role.ID == 0 {
		return errors.New("ID de rol no inválido")
	}

	role.Updated = time.Now()
	return s.repo.Update(role)
}

// DeleteRole usando el repository
func (s *RoleService) DeleteRole(id uint64) error {
	return s.repo.Delete(id)
}
