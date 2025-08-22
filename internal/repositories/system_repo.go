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

func (r *SystemRepository) GetPaginated(page, perPage int, nameQuery, descQuery string) ([]domain.System, int64, error) {
	var systems []domain.System
	var total int64

	query := r.db.Model(&domain.System{})

	if nameQuery != "" {
		query = query.Where("name LIKE ?", "%"+nameQuery+"%")
	}

	if descQuery != "" {
		query = query.Where("description LIKE ?", "%"+descQuery+"%")
	}

	// Contar el total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * perPage
	err := query.Offset(offset).Limit(perPage).Find(&systems).Error

	return systems, total, err
}

func (r *SystemRepository) GetByID(id uint64) (domain.System, error) {
	var system domain.System
	result := r.db.First(&system, id)
	if result.Error != nil {
		return domain.System{}, result.Error
	}
	return system, nil
}

func (r *SystemRepository) Create(system *domain.System) error {
	return r.db.Create(system).Error
}

func (r *SystemRepository) Update(system *domain.System) error {
	return r.db.Save(system).Error
}

func (r *SystemRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.System{}, id).Error
}
