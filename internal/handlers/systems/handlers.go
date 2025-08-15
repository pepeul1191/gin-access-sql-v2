package systems

import (
	"accessv2/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	service *services.SystemService
}

func NewSystemHandler(service *services.SystemService) *SystemHandler {
	return &SystemHandler{service: service}
}

func (h *SystemHandler) ListSystems(c *gin.Context) {
	systems, err := h.service.GetAllSystems()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "Error al obtener los sistemas",
		})
		return
	}

	globals, _ := c.Get("globals")
	c.HTML(http.StatusOK, "systems_list.html", gin.H{
		"title":   "Listado de Sistemas",
		"systems": systems,
		"globals": globals,
	})
}
