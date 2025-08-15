// config/bootstrap.go
package config

import (
	"accessv2/internal/handlers/auth"
	"accessv2/internal/handlers/common"
	"accessv2/internal/handlers/systems"
	"accessv2/internal/repositories"
	"accessv2/internal/services"
	"accessv2/pkg/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, store sessions.Store) *gin.Engine {
	router := gin.Default()

	// Configuración de cookies (seguridad)
	store.Options(sessions.Options{
		MaxAge:   86400 * 7, // 1 semana
		HttpOnly: true,      // Solo accesible por HTTP
		Secure:   true,      // Solo HTTPS en producción
		SameSite: http.SameSiteLaxMode,
	})

	// Middlewares globales (orden importante)
	router.Use(
		sessions.Sessions("mysession", store), // 1. Sesiones primero
		middleware.GlobalVarsMiddleware(),     // 2. Variables globales
		middleware.SessionMiddleware(),        // 3. Middleware de sesión personalizado
		middleware.CSRFMiddleware(),           // 4. CSRF (depende de las sesiones)
	)

	// Inicialización de repositorios
	systemRepo := repositories.NewSystemRepository(db)

	// Inicialización de servicios
	authService := services.NewAuthService()
	systemService := services.NewSystemService(systemRepo)

	// Inicialización de handlers
	commonHandler := common.NewCommonHandler()
	authHandler := auth.NewAuthHandler(authService)
	systemHandler := systems.NewSystemHandler(systemService)

	// Registrar rutas
	common.RegisterCommonRoutes(router, commonHandler)
	auth.RegisterAuthRoutes(router, authHandler)
	systems.RegisterSystemsRoutes(router, systemHandler)

	return router
}
