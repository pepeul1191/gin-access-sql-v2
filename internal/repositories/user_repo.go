// internal/repositories/user_repo.go
package repositories

import (
	"accessv2/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *UserRepository) GetPaginated(page, perPage int, usernameQuery, emailQuery string, statusQuery string) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Model(&domain.User{})

	if usernameQuery != "" {
		query = query.Where("username LIKE ?", "%"+usernameQuery+"%")
	}

	if emailQuery != "" {
		query = query.Where("email LIKE ?", "%"+emailQuery+"%")
	}

	// Filtro por estado
	if statusQuery != "" {
		if statusQuery == "active" {
			query = query.Where("activated = ?", true)
		} else if statusQuery == "inactive" {
			query = query.Where("activated = ?", false)
		}
	}

	// Contar el total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * perPage
	err := query.Offset(offset).Limit(perPage).Find(&users).Error

	return users, total, err
}

func (r *UserRepository) GetByID(id uint64) (domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.User{}, id).Error
}
