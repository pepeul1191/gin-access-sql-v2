package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/forms"
	"accessv2/internal/repositories"
	"accessv2/internal/responses"

	"accessv2/pkg/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
	db   *gorm.DB
}

func NewUserService(db *gorm.DB, repo *repositories.UserRepository) *UserService {
	return &UserService{
		db:   db,
		repo: repo}
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetPaginatedUsers(page, perPage int, usernameQuery, emailQuery string, statusFilter string) ([]domain.User, int64, error) {
	// Validación básica
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Delegar al repositorio
	return s.repo.GetPaginated(page, perPage, usernameQuery, emailQuery, statusFilter)
}

func (s *UserService) CreateUser(input *forms.UserCreateInput) (*domain.User, error) {

	var (
		ErrUserNameRequired = errors.New("el nombre del sistema es requerido")
	)
	// Validación de datos
	if input.Username == "" {
		return nil, ErrUserNameRequired
	}

	// Validación única
	err := s.repo.CheckUserExists(input.Username, input.Email, 0)
	if err != nil {
		return nil, errors.New("Usuario y/o correo en uso")
	}

	activated := false
	if input.Status == "active" {
		activated = true
	}

	// Crear objeto del dominio
	user := &domain.User{
		Username:      input.Username,
		Password:      input.Password,
		Email:         input.Email,
		ResetKey:      utils.RandomString(30),
		ActivationKey: utils.RandomString(30),
		Activated:     activated,
	}

	// Establecer fechas por defecto si no vienen
	if user.Created.IsZero() {
		user.Created = time.Now()
	}
	if user.Updated.IsZero() {
		user.Updated = user.Created
	}

	// Guardar en la base de datos
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// FetchSystem usando el repository (versión con pointer)
func (s *UserService) FetchUser(id uint64, user *domain.User) error {
	if user == nil {
		return errors.New("user pointer cannot be nil")
	}

	// Usar el método del repository que acepta pointer
	tempUser, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Copiar los valores al User proporcionado
	*user = tempUser
	return nil
}

// UpdateUser usando el repository
func (s *UserService) UpdateUser(user *domain.User) error {
	if user.ID == 0 {
		return errors.New("ID de usuario inválido")
	}

	err := s.repo.CheckUserExistsForUpdate(user.Username, user.Email, user.ID)
	if err != nil {
		return err // Si se encuentra un error (otro rol con el mismo nombre o correo), retornarlo.
	}

	user.Updated = time.Now()
	return s.repo.Update(user)
}

// DeleteSystem usando el repository
func (s *UserService) DeleteUser(id uint64) error {
	return s.repo.Delete(id)
}

func (s *UserService) ValidateBySystemUsernamePassword(systemID uint64, username, password string) (responses.UserWithAccess, error) {
	user, err := s.repo.GetBySystemUsernamePassword(systemID, username, password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.UserWithAccess{}, errors.New("Usuario y/o contraseña incorrectos")
		}
		return responses.UserWithAccess{}, fmt.Errorf("Error al validar usuario: %w", err)
	}

	if user.Activated == false {
		return responses.UserWithAccess{}, errors.New("Usuario no activo")
	}

	access, err := s.repo.GetUserNestedPermissionsBySystem(user.ID, systemID)
	if err != nil {
		return responses.UserWithAccess{}, err
	}

	userWithAccess := responses.UserWithAccess{
		User: responses.UserAccess{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			// Puedes incluir otros campos si los necesitas
		},
		Access: access,
	}

	return userWithAccess, nil
}
