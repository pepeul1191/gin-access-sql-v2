package users

import (
	"accessv2/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, handler *UserHandler) {

	usersGroup := r.Group("/users", middleware.AuthRequired())
	{
		usersGroup.GET("/", handler.ListUsers)
		usersGroup.POST("/create", handler.CreateUserHandler)
		usersGroup.GET("/create", handler.CreateUserHandler)
		usersGroup.POST("/:id/edit", handler.EditUserHandler)
		usersGroup.GET("/:id/edit", handler.EditUserHandler)
		usersGroup.GET("/:id/delete", handler.DeleteUserHandler)
	}
}
