// config/bootstrap.go
package config

import (
	"accessv2/internal/handlers/auth"
	"accessv2/internal/handlers/common"
	"accessv2/internal/services"
	"accessv2/pkg/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store sessions.Store) *gin.Engine {
	router := gin.Default()

	// 0. Configura el middlewares globales
	router.Use(sessions.Sessions("mysession", store))
	router.Use(middleware.GlobalVarsMiddleware())
	router.Use(middleware.GlobalVarsMiddleware())
	router.Use(middleware.SessionMiddleware())

	// 1. Inicializa los servicios
	authService := services.NewAuthService() // Crea la instancia del servicio

	// 2. Inicializa handlers
	commonHandler := common.NewCommonHandler()
	authHandler := auth.NewAuthHandler(authService)

	// Registrar rutas
	common.RegisterCommonRoutes(router, commonHandler)
	auth.RegisterAuthRoutes(router, authHandler)

	return router
}
