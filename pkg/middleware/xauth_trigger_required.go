// pkg/middleware/xauth_trigger_required.go
package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func XAuthTriggerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar si SECURE está habilitado
		secureEnabled, err := strconv.ParseBool(os.Getenv("SECURE"))
		if err != nil {
			// Si hay error al parsear, asumimos false por seguridad
			secureEnabled = false
		}

		// Si SECURE no está activado, continuamos sin validar
		if !secureEnabled {
			c.Next()
			return
		}

		// Validación original cuando SECURE=true
		incoming := c.GetHeader("X-Auth-Trigger")
		if incoming != os.Getenv("AUTH_HEADER") {
			fmt.Println("Unauthorized access attempt.")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or missing X-Auth-Trigger",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
