package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GlobalVars struct {
	AppName     string
	CurrentYear int
	Env         string // "development", "production", etc.
	BaseURL     string
	StaticURL   string
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
			BaseURL:     os.Getenv("BASE_URL"),
			StaticURL:   os.Getenv("STATIC_URL"),
			CurrentYear: time.Now().Year(),
		}
		// Colocar los bloques HTML en el contexto
		c.Set("globals", globals)
		// Continuar con la siguiente función del middleware o controlador
		c.Next()
	}
}
