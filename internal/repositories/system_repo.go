package repositories

import (
	"accessv2/internal/domain"
	"fmt"

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

	// Start with the base query and apply joins and filters first.
	query := r.db.Model(&domain.User{}).
		Joins("LEFT JOIN systems_users su ON users.id = su.user_id AND su.system_id = ?", systemID)

	// Apply filters
	if usernameQuery != "" {
		query = query.Where("users.username LIKE ?", "%"+usernameQuery+"%")
	}
	if emailQuery != "" {
		query = query.Where("users.email LIKE ?", "%"+emailQuery+"%")
	}
	if statusQuery != "" {
		if statusQuery == "active" {
			query = query.Where("users.activated = ?", true)
		} else if statusQuery == "inactive" {
			query = query.Where("users.activated = ?", false)
		}
	}

	// Count the total first, using a subquery to ensure accuracy with joins.
	var countQuery *gorm.DB
	countQuery = r.db.Table("(?) AS temp", query.Select("users.id")).Count(&total)
	if err := countQuery.Error; err != nil {
		return nil, 0, err
	}

	// Define the projection (columns to select)
	selects := "users.id, users.username, users.email, users.activated, CASE WHEN su.user_id IS NOT NULL THEN 1 ELSE 0 END AS association_status"

	// Finally, apply the custom SELECT and pagination before the Find call.
	offset := (page - 1) * perPage
	err := query.Select(selects).Offset(offset).Limit(perPage).Find(&users).Error

	fmt.Println("1 +++++++++++++++++++++++++++++++++++")
	fmt.Println(users)
	fmt.Println("2 +++++++++++++++++++++++++++++++++++")

	return users, total, err
}
