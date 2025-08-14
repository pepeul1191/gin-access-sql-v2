package common

import (
	"accessv2/pkg/middleware"
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
	sessionData, _ := c.Get("sessionData")

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":    "Página Principal",
		"globals":  globals,
		"nav_link": "",
		"session":  sessionData.(middleware.SessionData),
	})
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
