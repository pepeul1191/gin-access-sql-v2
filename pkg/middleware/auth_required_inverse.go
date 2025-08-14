// pkg/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequiredInverse redirige a la página principal si el usuario ya está autenticado
func AuthRequiredInverse() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		auth := session.Get("IsAuthenticated")

		// Si el usuario está autenticado, redirigir a la página principal
		if auth != nil {
			if isAuth, ok := auth.(bool); ok && isAuth {
				c.Redirect(http.StatusFound, "/")
				c.Abort()
				return
			}
		}

		// Si no está autenticado, continuar con la solicitud
		c.Next()
	}
}
