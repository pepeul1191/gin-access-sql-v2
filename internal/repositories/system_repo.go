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

func (r *SystemRepository) GetPaginatedUsers(page, perPage int, usernameQuery, emailQuery string, statusQuery string, systemID uint) ([]domain.UserSummary, int64, error) {
	var users []domain.UserSummary
	var total int64

	// Start with the base query and apply joins.
	query := r.db.Model(&domain.User{}).
		Joins("LEFT JOIN systems_users su ON users.id = su.user_id AND su.system_id = ?", systemID)

	// Apply filters. The logic here is key.
	if usernameQuery != "" {
		query = query.Where("users.username LIKE ?", "%"+usernameQuery+"%")
	}
	if emailQuery != "" {
		query = query.Where("users.email LIKE ?", "%"+emailQuery+"%")
	}

	// Filter for activated status only if the query parameter is provided.
	// The '2' option from your template handles all users, so we don't apply a filter for it.
	if statusQuery != "" {
		if statusQuery == "1" {
			// Filter for users who ARE in the system
			query = query.Where("su.user_id IS NOT NULL")
		} else if statusQuery == "0" {
			// Filter for users who ARE NOT in the system
			query = query.Where("su.user_id IS NULL")
		}
	}

	// Count the total first, using a subquery for accuracy with joins.
	// The count needs to reflect the filters, so it must be called on the built query.
	var countQuery *gorm.DB
	countQuery = r.db.Table("(?) AS temp", query.Select("users.id")).Count(&total)
	if err := countQuery.Error; err != nil {
		return nil, 0, err
	}

	// Finally, apply the custom SELECT and pagination before the Find call.
	// This order ensures your select and pagination are the final clauses in the query.
	selects := "users.id, users.username, users.email, users.activated, CASE WHEN su.user_id IS NOT NULL THEN 1 ELSE 0 END AS association_status"
	offset := (page - 1) * perPage

	err := query.Select(selects).Offset(offset).Limit(perPage).Find(&users).Error

	return users, total, err
}
