package services

import (
	"accessv2/internal/domain"
	"accessv2/internal/repositories"
	"errors"
	"time"

	"gorm.io/gorm"
)

// SystemUserService es la implementación del servicio.
type SystemUserService struct {
	repo *repositories.SystemUserRepository
	db   *gorm.DB
}

// NewSystemUserService crea una nueva instancia del servicio.
func NewSystemUserService(db *gorm.DB, repo *repositories.SystemUserRepository) *SystemUserService {
	return &SystemUserService{
		db:   db,
		repo: repo,
	}
}

// SaveSystemUsers se encarga de la lógica de negocio para crear o eliminar
// las relaciones entre usuarios y sistemas dentro de una transacción.
func (s *SystemUserService) SaveSystemUsers(systemID uint, items []domain.SystemUserItem) error {
	// 1. Iniciar la transacción.
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 2. Usar 'defer' para asegurar que la transacción siempre se revierta
	// si ocurre un error o si la función retorna antes del Commit.
	defer tx.Rollback()

	// 3. Iterar sobre los elementos recibidos.
	for _, item := range items {
		// Llamar al repositorio para buscar una relación existente dentro de la transacción.
		_, err := s.repo.FindSystemUser(tx, systemID, uint(item.ID))

		if item.Selected {
			// Caso 1: El usuario debe estar asociado al sistema.
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// Si la relación no existe, la creamos.
					newUser := &domain.SystemUser{
						UserID:   uint(item.ID),
						SystemID: systemID,
						Created:  time.Now(),
					}
					if err := s.repo.CreateSystemUser(tx, newUser); err != nil {
						return err
					}
				} else {
					// Si es otro tipo de error, lo retornamos.
					return err
				}
			}
			// Si la relación ya existe, no se hace nada.
		} else {
			// Caso 2: El usuario NO debe estar asociado al sistema.
			if err == nil {
				// Si la relación existe, la eliminamos.
				if err := s.repo.DeleteSystemUser(tx, systemID, uint(item.ID)); err != nil {
					return err
				}
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				// Si es un error de la base de datos (que no sea "registro no encontrado"), lo retornamos.
				return err
			}
			// Si la relación no existe, no se hace nada.
		}
	}

	// 4. Si el bucle finaliza sin errores, se confirma la transacción.
	return tx.Commit().Error
}
