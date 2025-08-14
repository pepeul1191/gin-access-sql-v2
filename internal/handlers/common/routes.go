package common

import (
	"accessv2/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCommonRoutes(r *gin.Engine, handler *CommonHandler) {
	// Rutas públicas
	r.GET("/", middleware.AuthRequired(), handler.Home)

	// Manejo de 404 (debe ser la última ruta)
	r.NoRoute(handler.NotFound)
}
