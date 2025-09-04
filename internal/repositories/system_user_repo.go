package repositories

import (
	"errors"

	"gorm.io/gorm"

	"accessv2/internal/domain"
)

// SystemUserRepository es la implementación del repositorio.
type SystemUserRepository struct {
	db *gorm.DB
}

// NewSystemUserRepository crea una nueva instancia del repositorio.
func NewSystemUserRepository(db *gorm.DB) *SystemUserRepository {
	return &SystemUserRepository{db: db}
}

// FindSystemUser busca una relación de usuario-sistema en la base de datos.
func (r *SystemUserRepository) FindSystemUser(tx *gorm.DB, systemID, userID uint) (*domain.SystemUser, error) {
	var user domain.SystemUser
	result := tx.Where("system_id = ? AND user_id = ?", systemID, userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// CreateSystemUser crea una nueva relación de usuario-sistema.
func (r *SystemUserRepository) CreateSystemUser(tx *gorm.DB, systemUser *domain.SystemUser) error {
	return tx.Create(systemUser).Error
}

// DeleteSystemUser elimina una relación de usuario-sistema.
func (r *SystemUserRepository) DeleteSystemUser(tx *gorm.DB, systemID, userID uint) error {
	return tx.Where("system_id = ? AND user_id = ?", systemID, userID).Delete(&domain.SystemUser{}).Error
}
