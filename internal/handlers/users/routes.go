package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, handler *UserHandler) {

	r.GET("/users", handler.ListUsers)
	r.POST("/users/create", handler.CreateUserHandler)
	r.GET("/users/create", handler.CreateUserHandler)
	r.POST("/users/:id/edit", handler.EditUserHandler)
	r.GET("/users/:id/edit", handler.EditUserHandler)
	r.GET("/users/:id/delete", handler.DeleteUserHandler)
}
