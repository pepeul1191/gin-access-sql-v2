// internal/routes/routes.go
package routes

import (
	"accessv2/internal/handlers"
	"accessv2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) *gin.Engine {
	// Logging
	r.Use(gin.Logger())

	// Ruta principal
	r.GET("/", handlers.Home)

	// Sign-in
	r.POST("/api/v1/sign-in", middlewares.SignInAuthenticate(), handlers.SignIn)

	// Grupo de rutas protegidas
	r.POST("/api/v1/files/:folder_name", middlewares.CheckJWT(), middlewares.FileValidation(), handlers.UploadFile)
	r.POST("/api/v1/public/:folder_name", middlewares.CheckJWT(), middlewares.FileValidation(), handlers.UploadFileToPublic)
	r.GET("/api/v1/files/:folder_name/:file_name", middlewares.CheckJWT(), handlers.DownloadFile)
	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Recurso no encontrado",
			"error":   c.Request.Method + " " + c.Request.URL.Path + " no existe",
		})
	})

	return r
}
