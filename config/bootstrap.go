package config

import (
	"accessv2/internal/handlers/auth"
	"accessv2/internal/handlers/common"
	"accessv2/internal/handlers/permissions"
	"accessv2/internal/handlers/roles"
	"accessv2/internal/handlers/systems"
	"accessv2/internal/handlers/users"
	"accessv2/internal/repositories"
	"accessv2/internal/services"

	"accessv2/pkg/middleware"
	"accessv2/pkg/utils"
	"html/template"
	"net/http"
	"time"

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

	// Cargar helpers a vistas
	router.SetFuncMap(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("02/01/2006 - 03:04:05 PM")
		},
		"add":     utils.Add,
		"sub":     utils.Sub,
		"scripts": utils.GenerateScriptsHTML,
		"styles":  utils.GenerateStylesHTML,
	})

	// Inicialización de repositorios
	systemRepo := repositories.NewSystemRepository(db)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	userSystemRepo := repositories.NewSystemUserRepository(db)
	userPermissionRepo := repositories.NewUserPermissionRepository(db)

	// Inicialización de servicios
	authService := services.NewAuthService()
	systemService := services.NewSystemService(systemRepo, userSystemRepo)
	permissionService := services.NewPermissionService(permissionRepo)
	userPermissionService := services.NewUserPermissionService(userPermissionRepo)
	userService := services.NewUserService(db, userRepo)
	roleService := services.NewRoleService(roleRepo)
	systemUserService := services.NewSystemUserService(db, userSystemRepo)

	// Inicialización de handlers
	commonHandler := common.NewCommonHandler()
	authHandler := auth.NewAuthHandler(authService)
	systemHandler := systems.NewSystemHandler(systemService, roleService, permissionService, systemUserService)
	userHandler := users.NewUserHandler(userService, userPermissionService)
	roleHandler := roles.NewRoleHandler(roleService)
	permissionHandler := permissions.NewPermissionHandler(permissionService)

	// Registrar rutas
	common.RegisterCommonRoutes(router, commonHandler)
	auth.RegisterAuthRoutes(router, authHandler)
	systems.RegisterSystemsRoutes(router, systemHandler, roleHandler, permissionHandler, userHandler)
	users.RegisterUserRoutes(router, userHandler)

	return router
}
