package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	// Dependencias si las necesitas (ej: templates, servicios)
}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) Home(c *gin.Context) {
	globals, _ := c.Get("globals")
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":   "Página Principal",
		"globals": globals,
	})
}

func (h *CommonHandler) SignIn(c *gin.Context) {
	globals, _ := c.Get("globals")
	c.HTML(http.StatusOK, "sign-in.html", gin.H{
		"title":   "Bienvenido",
		"globals": globals,
	})
}

func (h *CommonHandler) SignOut(c *gin.Context) {
	// Lógica de cierre de sesión
	c.Redirect(http.StatusFound, "/")
}

func (h *CommonHandler) NotFound(c *gin.Context) {
	if c.Request.Method == "GET" {
		globals, _ := c.Get("globals")
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title":   "Página no encontrada",
			"path":    c.Request.URL.Path,
			"globals": globals,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "not_found",
			"message": fmt.Sprintf("El recurso %s no existe", c.Request.URL.Path),
			"path":    c.Request.URL.Path,
		})
	}
}
