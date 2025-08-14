// pkg/middleware/auth.go
package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		auth := session.Get("IsAuthenticated")

		if auth == nil {
			log.Printf("Intento de acceso no autenticado a %s", c.Request.URL.Path)
			handleUnauthorized(c)
			return
		}

		if isAuth, ok := auth.(bool); !ok || !isAuth {
			log.Printf("Intento de acceso con sesión inválida desde %s", c.Request.RemoteAddr)
			handleUnauthorized(c)
			return
		}

		c.Next()
	}
}

func handleUnauthorized(c *gin.Context) {
	if isHTMLRequest(c) {
		globals, _ := c.Get("globals")
		c.HTML(http.StatusUnauthorized, "401.html", gin.H{
			"title":   "Acceso no autorizado",
			"globals": globals,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "authentication_required",
			"message": "Por favor inicie sesión",
		})
	}
	c.Abort()
}

func isHTMLRequest(c *gin.Context) bool {
	return c.Request.Method == "GET"
}
