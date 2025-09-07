package forms

type UserCreateInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password"`
	Email    string `form:"email"`
	Status   string `form:"status"`
}

type UserEditInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password"`
	Email    string `form:"email"`
	Status   string `form:"status"`
}

// LoginRequest representa la estructura del JSON de entrada
type SignInRequest struct {
	SystemID uint64 `json:"system_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
