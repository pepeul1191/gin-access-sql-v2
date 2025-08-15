package repositories

import (
	"accessv2/internal/domain"

	"gorm.io/gorm"
)

type SystemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{db: db}
}

func (r *SystemRepository) GetAll() ([]domain.System, error) {
	var systems []domain.System
	result := r.db.Find(&systems)
	if result.Error != nil {
		return nil, result.Error
	}
	return systems, nil
}
