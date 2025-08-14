package config

import (
	"accessv2/internal/handlers/common"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Handlers
	commonHandler := common.NewCommonHandler()
	//authHandler := auth.NewAuthHandler() // Asumiendo que tienes este handler

	// Registrar rutas de cada módulo
	common.RegisterCommonRoutes(router, commonHandler)
	//auth.RegisterAuthRoutes(router, authHandler) // Ejemplo para auth

	return router
}
