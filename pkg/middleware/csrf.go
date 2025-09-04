// pkg/middleware/csrf.go
package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const csrfTokenLength = 32

// generateToken crea un token CSRF seguro
func generateToken() string {
	b := make([]byte, csrfTokenLength)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// isAPIRequest determina si es una ruta API
func isAPIRequest(c *gin.Context) bool {
	return strings.HasPrefix(c.Request.URL.Path, "/api/")
}

// isAjaxRequest detecta peticiones AJAX/HTTP modernas
func isAjaxRequest(c *gin.Context) bool {
	return c.GetHeader("X-Requested-With") == "XMLHttpRequest" ||
		strings.Contains(c.GetHeader("Accept"), "application/json") ||
		c.GetHeader("Content-Type") == "application/json"
}

// getSubmittedToken obtiene el token según el tipo de petición
func getSubmittedToken(c *gin.Context) string {
	// 1. Buscar en headers (para AJAX/APIs)
	if token := c.GetHeader("X-CSRFToken"); token != "" {
		return token
	}

	// 2. Buscar en form-data (para formularios tradicionales)
	if token := c.PostForm("_csrf"); token != "" {
		return token
	}

	// 3. Buscar en JSON body (para APIs que usan JSON)
	if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		var jsonBody struct {
			CSRFToken string `json:"csrf_token"`
		}
		if err := c.ShouldBindJSON(&jsonBody); err == nil && jsonBody.CSRFToken != "" {
			return jsonBody.CSRFToken
		}
	}

	return ""
}

// CSRFMiddleware implementa protección CSRF inteligente
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Excluir rutas API
		if isAPIRequest(c) {
			c.Next()
			return
		}

		session := sessions.Default(c)

		// Generar token para métodos seguros (GET, HEAD, OPTIONS)
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {

			token := generateToken()
			session.Set("csrf_token", token)
			if err := session.Save(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "No se pudo guardar el token CSRF",
				})
				return
			}
			c.Set("csrf_token", token) // Para usar en templates
			c.Next()
			return
		}

		// Validar para métodos peligrosos (POST, PUT, PATCH, DELETE)
		storedToken, ok := session.Get("csrf_token").(string)
		if !ok || storedToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Token CSRF no encontrado",
				"message": "La sesión no contiene un token CSRF válido",
			})
			return
		}

		globals, _ := c.Get("globals")

		// Obtener token enviado
		submittedToken := getSubmittedToken(c)
		if submittedToken == "" {
			if isAjaxRequest(c) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":   "Token CSRF faltante",
					"message": "Debe incluir el token CSRF en headers o formulario",
				})
			} else {
				c.HTML(http.StatusForbidden, "403", gin.H{
					"title":    "Error de seguridad",
					"message":  "El formulario no contenía el token de seguridad requerido",
					"code":     403,
					"globals":  globals,
					"styles":   []string{"css/common"},
					"scripts":  []string{"js/403"},
					"solution": "Recargue la página e intente nuevamente",
				})
			}
			return
		}

		// Comparación segura contra timing attacks
		if subtle.ConstantTimeCompare([]byte(storedToken), []byte(submittedToken)) != 1 {
			if isAjaxRequest(c) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":   "Token CSRF inválido",
					"message": "El token proporcionado no coincide",
				})
			} else {
				c.HTML(http.StatusForbidden, "403.html", gin.H{
					"title":    "Error de seguridad",
					"message":  "El token proporcionado no coincide",
					"code":     403,
					"globals":  globals,
					"solution": "Recargue la página e intente nuevamente",
				})
			}
			return
		}

		c.Next()
	}
}
