package common

import "github.com/gin-gonic/gin"

func RegisterCommonRoutes(r *gin.Engine, handler *CommonHandler) {
	// Rutas públicas
	r.GET("/", handler.Home)
	r.GET("/sign-in", handler.SignIn)
	r.GET("/sign-out", handler.SignOut)

	// Manejo de 404 (debe ser la última ruta)
	r.NoRoute(handler.NotFound)
}
