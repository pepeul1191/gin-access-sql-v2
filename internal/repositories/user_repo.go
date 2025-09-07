// internal/repositories/user_repo.go
package repositories

import (
	"accessv2/internal/domain"
	"accessv2/internal/responses"
	"errors"

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

func (r *UserRepository) CheckUserExists(username, email string, excludeID uint) error {
	var existingUser domain.User
	query := r.db.Model(&domain.User{}).
		Where("username = ? OR email = ?", username, email)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	result := query.First(&existingUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // No existe, todo bien
		}
		return result.Error // Error de base de datos
	}

	// Determinar qué campo causó el conflicto
	if existingUser.Username == username {
		return errors.New("username already exists")
	}
	if existingUser.Email == email {
		return errors.New("email already exists")
	}

	return nil
}

func (r *UserRepository) CheckUserExistsForUpdate(username string, email string, id uint) error {
	var existingUser domain.User

	// La consulta busca un usuario cuyo 'username' o 'email' coincida con los valores proporcionados,
	// pero que su 'id' sea diferente al del usuario actual.
	query := r.db.Model(&domain.User{}).
		Where("(username = ? OR email = ?) AND id != ?", username, email, id)

	result := query.First(&existingUser)

	// Si no se encuentra ningún registro, significa que los datos no están en uso por otro usuario.
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil // No se encontró un usuario con ese nombre o correo, todo bien.
		}
		// Se encontró un error de base de datos.
		return result.Error
	}

	// Si el resultado no es un error, GORM encontró un registro.
	// Esto significa que ya existe un usuario con el mismo nombre de usuario o correo.
	if existingUser.Username == username {
		return errors.New("El nombre de usuario ya está en uso por otro usuario.")
	}
	if existingUser.Email == email {
		return errors.New("El correo electrónico ya está en uso por otro usuario.")
	}

	return errors.New("El nombre de usuario o correo electrónico ya están en uso.")
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

func (r *UserRepository) GetBySystemUsernamePassword(systemID uint64, username, password string) (domain.User, error) {
	var user domain.User

	result := r.db.
		Joins("JOIN systems_users ON systems_users.user_id = users.id").
		Where("systems_users.system_id = ? AND users.username = ? AND users.password = ?", systemID, username, password).
		First(&user)

	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user, nil
}

func (r *UserRepository) GetUserNestedPermissionsBySystem(userID uint, systemID uint64) (responses.SystemAccess, error) {
	var flatPermissions []domain.UserSystemPermission

	query := `
        SELECT
            S.id AS system_id,
            S.name AS system_name,
            R.id AS role_id,
            R.name AS role_name,
            P.id AS permission_id,
            P.name AS permission_name
        FROM systems_users_permissions AS SUP
        INNER JOIN users AS U ON SUP.user_id = U.id
        INNER JOIN permissions AS P ON SUP.permission_id = P.id
        INNER JOIN roles AS R ON P.role_id = R.id
        INNER JOIN systems AS S ON R.system_id = S.id
        WHERE SUP.user_id = ? AND S.id = ?;
    `

	if err := r.db.Raw(query, userID, systemID).Scan(&flatPermissions).Error; err != nil {
		return responses.SystemAccess{}, err
	}

	// If no permissions found, return empty system access
	if len(flatPermissions) == 0 {
		return responses.SystemAccess{Roles: []*responses.RoleAccess{}}, nil
	}

	// Reconstruct the nested structure with only necessary fields
	rolesMap := make(map[uint64]*responses.RoleAccess)
	var roles []*responses.RoleAccess

	for _, p := range flatPermissions {
		// Find or create the role
		role, roleExists := rolesMap[p.RoleID]
		if !roleExists {
			role = &responses.RoleAccess{
				ID:          uint(p.RoleID),
				Name:        p.RoleName,
				Permissions: []responses.PermissionAccess{},
			}
			rolesMap[p.RoleID] = role
			roles = append(roles, role)
		}

		// Add the permission to the role
		role.Permissions = append(role.Permissions, responses.PermissionAccess{
			ID:   uint(p.PermissionID),
			Name: p.PermissionName,
		})
	}

	// Create the simplified system access response
	systemAccess := responses.SystemAccess{
		Roles: roles,
	}

	return systemAccess, nil
}
