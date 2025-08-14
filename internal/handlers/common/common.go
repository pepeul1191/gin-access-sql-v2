package common

import (
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
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Página Principal",
	})
}

func (h *CommonHandler) SignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "sign-in.html", gin.H{})
}

func (h *CommonHandler) SignOut(c *gin.Context) {
	// Lógica de cierre de sesión
	c.Redirect(http.StatusFound, "/")
}

func (h *CommonHandler) NotFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", gin.H{
		"title": "Página no encontrada",
	})
}
