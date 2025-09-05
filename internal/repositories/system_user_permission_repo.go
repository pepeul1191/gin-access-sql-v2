package repositories

import (
	"accessv2/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type UserPermissionRepository struct {
	db *gorm.DB
}

// Nuevo repositorio para el acceso a datos
func NewUserPermissionRepository(db *gorm.DB) *UserPermissionRepository {
	return &UserPermissionRepository{db: db}
}

// Insertar permisos si no existen
func (r *UserPermissionRepository) InsertPermissions(permissions []domain.SystemUserPermission) error {
	for _, perm := range permissions {
		// Verificar si el permiso ya existe en la base de datos
		var existingPermission domain.SystemUserPermission
		if err := r.db.Where("system_id = ? AND user_id = ? AND permission_id = ?", perm.SystemID, perm.UserID, perm.PermissionID).First(&existingPermission).Error; err == nil {
			// Si ya existe, no lo insertamos
			continue
		}
		// Si no existe, insertar el nuevo permiso
		if err := r.db.Create(&perm).Error; err != nil {
			return err
		}
	}
	return nil
}

// Eliminar permisos previos para un usuario en el sistema que no estén en la lista
func (r *UserPermissionRepository) DeletePermissions(systemID, userID uint, permissionIDs []uint64) error {
	// Eliminar permisos que no están en la lista de permisos seleccionados
	if err := r.db.Where("system_id = ? AND user_id = ? AND permission_id NOT IN (?)", systemID, userID, permissionIDs).Delete(&domain.SystemUserPermission{}).Error; err != nil {
		return err
	}
	return nil
}
func (r *UserPermissionRepository) GetSystemUserRolesPermissions(systemID, userID uint64) ([]domain.SystemUserRolesPermissions, error) {

	var permissions []domain.SystemUserRolesPermissions

	query := `
        SELECT
            su.user_id,
            su.system_id,
            p.id AS permission_id,
            p.name AS permission_name,
            r.id AS role_id,
            r.name AS role_name,
            CASE
                WHEN sup.id IS NOT NULL THEN 1
                ELSE 0
            END AS is_assigned
        FROM systems_users su
        JOIN roles r ON r.system_id = su.system_id
        JOIN permissions p ON p.role_id = r.id
        LEFT JOIN systems_users_permissions sup
            ON sup.system_id = su.system_id
            AND sup.user_id = su.user_id
            AND sup.permission_id = p.id
        WHERE su.system_id = ? AND su.user_id = ?;
    `

	result := r.db.Raw(query, systemID, userID).Scan(&permissions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []domain.SystemUserRolesPermissions{}, nil
		}
		return nil, result.Error
	}

	return permissions, nil
}
