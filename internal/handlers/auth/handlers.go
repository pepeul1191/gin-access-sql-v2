// internal/handlers/auth/handlers.go
package auth

import (
	"net/http"

	"accessv2/internal/forms"
	"accessv2/internal/services"
	"accessv2/pkg/middleware"
	"accessv2/pkg/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	globals, _ := c.Get("globals")
	session := sessions.Default(c)

	if c.Request.Method == "GET" {
		// Recuperar flashes al mostrar el formulario
		flashes := session.Flashes("error")
		session.Save()

		csrfToken, _ := c.Get("csrf_token")

		c.HTML(http.StatusOK, "sign-in", gin.H{
			"title":       "Iniciar Sesión",
			"globals":     globals,
			"csrfToken":   csrfToken,
			"flash_error": utils.FirstFlashOrEmpty(flashes),
			"styles":      []string{"css/auth"},
			"scripts":     []string{},
		})
		return
	}

	// Procesar POST
	var form forms.LoginForm
	if err := c.ShouldBind(&form); err != nil {
		session.AddFlash("Por favor completa todos los campos", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/sign-in") // Redirige en lugar de renderizar
		return
	}

	isValid, err := h.authService.Authenticate(form.Username, form.Password)
	if err != nil || !isValid {
		session.AddFlash("Usuario o contraseña incorrectos", "error")
		session.Save()
		c.Redirect(http.StatusFound, "/sign-in") // Redirige en lugar de renderizar
		return
	}

	// Login exitoso
	session.Set("IsAuthenticated", true)
	session.Set("Username", form.Username)
	session.Set("UserID", "1")
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	session := sessions.Default(c)

	// 1. Destruir la sesión (múltiples métodos)
	session.Clear()                               // Elimina todos los valores
	session.Options(sessions.Options{MaxAge: -1}) // Expira la cookie
	session.Save()                                // Guardar cambios

	// 2. Mostrar template de confirmación
	c.HTML(http.StatusOK, "sign-out", gin.H{
		"title":   "Sesión cerrada",
		"styles":  []string{"css/common"},
		"globals": c.MustGet("globals"),
	})
}

func (h *AuthHandler) Session(c *gin.Context) {
	// 1. Obtener los datos de sesión del contexto
	sessionData, exists := c.Get("sessionData")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se encontraron datos de sesión",
		})
		return
	}

	// 2. Hacer type assertion para obtener el struct SessionData
	session, ok := sessionData.(middleware.SessionData)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Formato de sesión inválido",
		})
		return
	}

	// 3. Devolver los datos de sesión en JSON
	c.JSON(http.StatusOK, gin.H{
		"is_authenticated": session.IsAuthenticated,
		"username":         session.Username,
		"user_id":          session.UserID,
		// ... otros campos si los necesitas
	})
}
