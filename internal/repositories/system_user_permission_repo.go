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

func (r *UserPermissionRepository) DeletePermissions(systemID, userID, roleID uint) error {
	// Step 1: Create a subquery to select all permission IDs belonging to the specified role.
	// We use `Model` to specify the table and `Select("id")` to get only the IDs.
	subQuery := r.db.Model(&domain.Permission{}).Select("id").Where("role_id = ?", roleID)

	// Step 2: Delete records from the systems_users_permissions table.
	// We use the subquery in the `IN` clause to target only the permissions from the specified role.
	// The `system_id` and `user_id` are included for safety and precision.
	if err := r.db.Where(
		"system_id = ? AND user_id = ? AND permission_id IN (?)",
		systemID,
		userID,
		subQuery,
	).Delete(&domain.SystemUserPermission{}).Error; err != nil {
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

// GetUserNestedPermissions fetches and structures a user's permissions hierarchically.
func (r *UserPermissionRepository) GetUserNestedPermissions(userID uint) ([]domain.System, error) {
	var flatPermissions []domain.UserSystemPermission

	query := `
		SELECT
				S.id AS system_id,
				S.name AS system_name,
				R.id AS role_id,
				R.name AS role_name,
				P.id AS permission_id,
				P.name AS permission_name
		FROM systems_users_permissions AS SUP
		INNER JOIN users AS U ON SUP.user_id = U.id
		INNER JOIN permissions AS P ON SUP.permission_id = P.id
		INNER JOIN roles AS R ON P.role_id = R.id
		INNER JOIN systems AS S ON R.system_id = S.id
		WHERE SUP.user_id = ?;
    `

	if err := r.db.Raw(query, userID).Scan(&flatPermissions).Error; err != nil {
		return nil, err
	}

	// Now, reconstruct the nested structure in Go.
	systemsMap := make(map[uint64]*domain.System)
	rolesMap := make(map[uint64]*domain.Role)

	for _, p := range flatPermissions {
		// Find or create the system.
		system, systemExists := systemsMap[p.SystemID]
		if !systemExists {
			system = &domain.System{
				ID:    uint(p.SystemID),
				Name:  p.SystemName,
				Roles: []*domain.Role{}, // Slice of pointers
			}
			systemsMap[p.SystemID] = system
		}

		// Find or create the role within the system.
		role, roleExists := rolesMap[p.RoleID]
		if !roleExists {
			role = &domain.Role{
				ID:          uint(p.RoleID),
				Name:        p.RoleName,
				Permissions: []domain.Permission{},
			}
			rolesMap[p.RoleID] = role
			system.Roles = append(system.Roles, role) // Append the pointer directly
		}

		// Add the permission to the role.
		role.Permissions = append(role.Permissions, domain.Permission{
			ID:   uint(p.PermissionID),
			Name: p.PermissionName,
		})
	}

	// Convert the map to a slice for the final output.
	var result []domain.System
	for _, s := range systemsMap {
		result = append(result, *s)
	}

	return result, nil
}
