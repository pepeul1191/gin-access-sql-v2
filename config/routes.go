// config/routes.go
package config

import (
	"accessv2/internal/handlers/common"
	"accessv2/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func SetupRouter(store sessions.Store) *gin.Engine {
	router := gin.Default()
	// Middleware de variables globales
	router.Use(middleware.GlobalVarsMiddleware())
	// Handlers
	commonHandler := common.NewCommonHandler()
	//authHandler := auth.NewAuthHandler() // Asumiendo que tienes este handler

	// Registrar rutas de cada módulo
	common.RegisterCommonRoutes(router, commonHandler, store)
	//auth.RegisterAuthRoutes(router, authHandler) // Ejemplo para auth

	return router
}
