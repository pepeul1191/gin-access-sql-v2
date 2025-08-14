package auth

import (
	"accessv2/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, handler *AuthHandler) {
	r.GET("/sign-in", middleware.AuthRequiredInverse(), handler.SignIn)
	r.POST("/sign-in", middleware.AuthRequiredInverse(), handler.SignIn)
	r.GET("/sign-out", handler.SignOut)
	r.GET("/session", middleware.AuthRequired(), handler.Session)
}
