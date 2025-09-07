// internal/responses/user_responses.go
package responses

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
	User   UserAccess   `json:"user"`
	Access SystemAccess `json:"access"`
}

type UserAccess struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Solo incluir los campos que necesitas mostrar
}

type SignResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    UserWithAccess `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}
