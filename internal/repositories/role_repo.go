package repositories

import (
	"accessv2/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CheckRoleExistsInSystem(name string, systemID int) error {
	var existingRole domain.Role
	query := r.db.Model(&domain.Role{}).
		Where("name = ? OR system_id = ?", name, systemID)

	result := query.First(&existingRole)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // No existe, todo bien
		}
		return result.Error // Error de base de datos
	}

	return errors.New("Nombre de usuario ya en uso en el sistema")
}

func (r *RoleRepository) GetRolesBySystemID(systemID int) ([]domain.Role, error) {
	var roles []domain.Role

	query := r.db.Model(&domain.Role{}).Where("system_id = ?", systemID)

	err := query.Find(&roles).Error

	return roles, err
}

func (r *RoleRepository) GetByID(id uint64) (domain.Role, error) {
	var role domain.Role
	result := r.db.First(&role, id)
	if result.Error != nil {
		return domain.Role{}, result.Error
	}
	return role, nil
}

func (r *RoleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *domain.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Role{}, id).Error
}
