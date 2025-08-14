// pkg/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Versión que acepta sesiones opcionales
func AuthRequired(store ...sessions.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Si se pasó un session store
		if len(store) > 0 {
			session, _ := store[0].Get(c.Request, "session-name")
			if auth, ok := session.Values["authenticated"].(bool); ok && auth {
				c.Next()
				return
			}
		}

		// Manejar no autenticado
		if c.Request.Method == "GET" {
			globals, _ := c.Get("globals")
			c.HTML(http.StatusUnauthorized, "401.html", gin.H{
				"title":   "Acceso no autorizado",
				"globals": globals,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
		}
		c.Abort()
	}
}
