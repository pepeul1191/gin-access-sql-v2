package repositories

import (
	"accessv2/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) CheckPermissionExistsInRole(name string, roleID int) error {
	var existingPermission domain.Permission
	query := r.db.Model(&domain.Permission{}).
		Where("name = ? AND role_id = ?", name, roleID)

	result := query.First(&existingPermission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // No existe, todo bien
		}
		return result.Error // Error de base de datos
	}

	return errors.New("Nombre de permiso ya en uso en el rol")
}

func (r *PermissionRepository) CheckPermissionNameExistsForUpdate(name string, roleID int, permissionID int) error {
	var existingPermission domain.Permission
	query := r.db.Model(&domain.Permission{}).
		Where("name = ? AND role_id = ? AND id != ?", name, roleID, permissionID)

	result := query.First(&existingPermission)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // No existe otro rol con ese nombre, todo bien
		}
		return result.Error // Error de base de datos
	}

	// Se encontró un rol con el mismo nombre y roleID, pero con un ID diferente
	return errors.New("Ya existe un permiso con este nombre en el rol")
}

func (r *PermissionRepository) GetPaginated(page, perPage int, roleID int) ([]domain.Permission, int64, error) {
	var permissions []domain.Permission
	var total int64

	query := r.db.Model(&domain.Permission{}).Where("role_id = ?", roleID)

	// Contar el total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * perPage
	err := query.Offset(offset).Limit(perPage).Find(&permissions).Error

	return permissions, total, err
}

func (r *PermissionRepository) GetPermissionsByRoleID(roleID int) ([]domain.Permission, error) {
	var permissions []domain.Permission

	query := r.db.Model(&domain.Permission{}).Where("role_id = ?", roleID)

	err := query.Find(&permissions).Error

	return permissions, err
}

func (r *PermissionRepository) GetByID(id uint64) (domain.Permission, error) {
	var permission domain.Permission
	result := r.db.First(&permission, id)
	if result.Error != nil {
		return domain.Permission{}, result.Error
	}
	return permission, nil
}

func (r *PermissionRepository) Create(permission *domain.Permission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) Update(permission *domain.Permission) error {
	return r.db.Save(permission).Error
}

func (r *PermissionRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Permission{}, id).Error
}
