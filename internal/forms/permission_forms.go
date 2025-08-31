package forms

type PermissionCreateInput struct {
	Name string `form:"name" binding:"required"`
}

type PermissionEditInput struct {
	Name string `form:"name" binding:"required"`
}
