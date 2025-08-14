// internal/middlewares/cors.go
package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si CORS está habilitado
		corsEnabled := os.Getenv("CORS_ENABLED") == "true"
		if !corsEnabled {
			c.Next()
			return
		}

		// Obtener hosts permitidos desde .env
		allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
		allowedMethods := os.Getenv("ALLOWED_METHODS")
		allowedHeaders := os.Getenv("ALLOWED_HEADERS")
		// Verificar si el origen está permitido
		requestOrigin := c.Request.Header.Get("Origin")
		originAllowed := false

		// Si es *, permitir cualquier origen (no recomendado para producción)
		if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
			originAllowed = true
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Verificar contra la lista de orígenes permitidos
			for _, origin := range allowedOrigins {
				if origin == requestOrigin {
					originAllowed = true
					c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
					break
				}
			}
		}

		// Configurar headers CORS
		if originAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 horas
		}

		// Manejar solicitudes OPTIONS
		if c.Request.Method == "OPTIONS" {
			if originAllowed {
				c.AbortWithStatus(http.StatusNoContent)
				return
			} else {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Next()
	}
}
