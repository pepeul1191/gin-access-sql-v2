package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/forms"
	"accessv2/internal/repositories"
	"errors"
	"time"
)

type PermissionService struct {
	repo *repositories.PermissionRepository
}

func NewPermissionService(repo *repositories.PermissionRepository) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) GetPaginatedRolePermissions(page, perPage int, roleID int) ([]domain.Permission, int64, error) {
	// Validación básica
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Delegar al repositorio
	return s.repo.GetPaginated(page, perPage, roleID)
}

func (s *PermissionService) GetAllByRoleID(roleID int) ([]domain.Permission, error) {
	// Delegar al repositorio
	return s.repo.GetPermissionsByRoleID(roleID)
}

func (s *PermissionService) CreatePermission(input *forms.PermissionCreateInput, roleID int) (*domain.Permission, error) {

	var (
		ErrNameRequired = errors.New("El nombre del permiso es requerido")
	)
	// Validación de datos
	if input.Name == "" {
		return nil, ErrNameRequired
	}

	// Validación única
	err := s.repo.CheckPermissionExistsInRole(input.Name, roleID)
	if err != nil {
		return nil, errors.New("Nombre de permiso ya en uso en el rol")
	}
	// Crear objeto del dominio
	permission := &domain.Permission{
		Name:   input.Name,
		RoleID: uint(roleID),
	}

	// Establecer fechas por defecto si no vienen
	if permission.Created.IsZero() {
		permission.Created = time.Now()
	}
	if permission.Updated.IsZero() {
		permission.Updated = permission.Created
	}

	// Guardar en la base de datos
	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// FetchSystem usando el repository (versión con pointer)
func (s *PermissionService) FetchPermission(id uint64, permission *domain.Permission) error {
	if permission == nil {
		return errors.New("Permission pointer cannot be nil")
	}

	// Usar el método del repository que acepta pointer
	temp, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Copiar los valores al User proporcionado
	*permission = temp
	return nil
}

// UpdateRole usando el repository
func (s *PermissionService) UpdatePermssion(permission *domain.Permission) error {
	if permission.ID == 0 {
		return errors.New("ID de rol no inválido")
	}

	err := s.repo.CheckPermissionNameExistsForUpdate(permission.Name, int(permission.RoleID), int(permission.ID))
	if err != nil {
		return err // Si se encuentra un error (otro rol con el mismo nombre), retornarlo.
	}

	permission.Updated = time.Now()
	return s.repo.Update(permission)
}

// DeleteRole usando el repository
func (s *PermissionService) DeletePermission(id uint64) error {
	return s.repo.Delete(id)
}
