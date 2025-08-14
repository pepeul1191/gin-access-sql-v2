// pkg/middleware/global_vars
package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GlobalVars struct {
	AppName           string
	CurrentYear       int
	Env               string // "development", "production", etc.
	GoogleAnalyticsID string
	// Agrega aquí todas las variables globales que necesites
}

func GlobalVarsMiddleware() gin.HandlerFunc {
	// Intenta cargar el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("⚠️ No se pudo cargar el archivo .env: %v", err)
		log.Println("⚠️ Usando variables de entorno del sistema")
	}

	return func(c *gin.Context) {
		// Inyecta variables globales desde .env o sistema
		globals := GlobalVars{
			AppName:     os.Getenv("APP_NAME"),
			CurrentYear: time.Now().Year(),
		}

		c.Set("globals", globals)

		c.Next()
	}
}
