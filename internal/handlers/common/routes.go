package common

import (
	"accessv2/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func RegisterCommonRoutes(r *gin.Engine, handler *CommonHandler, store sessions.Store) {
	// Rutas públicas
	r.GET("/", middleware.AuthRequired(store), handler.Home)
	r.GET("/sign-in", handler.SignIn)
	r.GET("/sign-out", handler.SignOut)

	// Manejo de 404 (debe ser la última ruta)
	r.NoRoute(handler.NotFound)
}
