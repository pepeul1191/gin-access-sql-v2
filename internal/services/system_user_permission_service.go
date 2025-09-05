package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/repositories"
	"errors"
	"time"
)

type UserPermissionService struct {
	repo *repositories.UserPermissionRepository
}

// Crear un nuevo servicio
func NewUserPermissionService(repo *repositories.UserPermissionRepository) *UserPermissionService {
	return &UserPermissionService{repo: repo}
}

func (s *UserPermissionService) GetUserRolesAndPermissions(systemID uint64, userID uint64) ([]domain.RoleWithPermissions, error) {
	if s.repo == nil {
		return nil, errors.New("repository is not initialized")
	}
	flatPermissions, err := s.repo.GetSystemUserRolesPermissions(systemID, userID)
	if err != nil {
		return nil, err
	}
	if flatPermissions == nil {
		return nil, errors.New("no roles or permissions found")
	}

	rolesMap := make(map[uint64]*domain.RoleWithPermissions)

	for _, p := range flatPermissions {
		if _, exists := rolesMap[uint64(p.RoleID)]; !exists {
			rolesMap[uint64(p.RoleID)] = &domain.RoleWithPermissions{
				ID:   p.RoleID,
				Name: p.RoleName,
			}
		}

		role := rolesMap[uint64(p.RoleID)]
		role.Permissions = append(role.Permissions, domain.UserPermission{
			ID:         p.PermissionID,
			Name:       p.PermissionName,
			IsAssigned: p.IsAssigned,
		})
	}

	var result []domain.RoleWithPermissions
	for _, role := range rolesMap {
		result = append(result, *role)
	}

	return result, nil
}

func (s *UserPermissionService) AssociatePermissions(systemID uint, userID uint, roleID uint, permissionIDs []uint64) error {
	// Eliminar los permisos previos que no están en la lista de permisos seleccionados
	if err := s.repo.DeletePermissions(systemID, userID, roleID); err != nil {
		return err
	}

	// Crear los registros de permisos que serán insertados
	var permissions []domain.SystemUserPermission
	for _, permID := range permissionIDs {
		permissions = append(permissions, domain.SystemUserPermission{
			SystemID:     systemID,
			UserID:       userID,
			PermissionID: uint(permID),
			Created:      time.Now(),
		})
	}

	// Insertar los nuevos permisos
	if err := s.repo.InsertPermissions(permissions); err != nil {
		return err
	}

	return nil
}
