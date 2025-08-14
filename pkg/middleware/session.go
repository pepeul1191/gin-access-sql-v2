// pkg/middleware/session.go
package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionData representa los datos de la sesión que queremos exponer.
type SessionData struct {
	IsAuthenticated bool
	Username        string
	UserID          int
	// ... otros campos
}

// SessionMiddleware extrae los datos de la sesión y los guarda en el contexto de Gin.
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		data := SessionData{
			IsAuthenticated: session.Get("IsAuthenticated") != nil,
			Username:        getStringFromSession(session, "Username"),
			UserID:          getIntFromSession(session, "UserID"),
		}

		// Guardamos el struct en el contexto de Gin
		c.Set("sessionData", data)

		c.Next() // Continuar con los demás middlewares/handlers
	}
}

// --- Helpers para manejo seguro de tipos ---

func getStringFromSession(session sessions.Session, key string) string {
	if val := session.Get(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return "" // Valor por defecto
}

func getIntFromSession(session sessions.Session, key string) int {
	if val := session.Get(key); val != nil {
		if num, ok := val.(int); ok {
			return num
		}
	}
	return 0 // Valor por defecto
}
