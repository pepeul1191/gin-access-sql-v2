// internal/responses/user_responses.go
package responses

import (
	"github.com/golang-jwt/jwt/v5"
)

// En tu archivo domain/user.go o donde tengas las estructuras
type SystemAccess struct {
	Roles []*RoleAccess `json:"roles"`
}

type RoleAccess struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Permissions []PermissionAccess `json:"permissions"`
}

type PermissionAccess struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserWithAccess struct {
	User  UserAccess    `json:"user"`
	Roles []*RoleAccess `json:"roles"`
}

type UserAccess struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
	// Solo incluir los campos que necesitas mostrar
}

type SignResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    UserWithAccess `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

type CustomClaims struct {
	UserID   uint64        `json:"user_id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Roles    []*RoleAccess `json:"roles"`
	jwt.RegisteredClaims
}
