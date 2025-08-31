package forms

type RoleCreateInput struct {
	Name string `form:"name" binding:"required"`
}

type RoleEditInput struct {
	Name string `form:"name" binding:"required"`
}
