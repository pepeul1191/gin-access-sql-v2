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
