package systems

import (
	"accessv2/internal/services"
	"accessv2/pkg/middleware"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	service *services.SystemService
}

func NewSystemHandler(service *services.SystemService) *SystemHandler {
	return &SystemHandler{service: service}
}

func (h *SystemHandler) ListSystems(c *gin.Context) {
	// Obtener parámetros de paginación y búsqueda
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	nameQuery := strings.TrimSpace(c.Query("name"))
	descQuery := strings.TrimSpace(c.Query("description"))

	// Validar parámetros
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Obtener sistemas paginados
	systems, total, err := h.service.GetPaginatedSystems(page, perPage, nameQuery, descQuery)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "Error al obtener los sistemas",
		})
		return
	}

	// Calcular total de páginas
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	// Calcular registros mostrados
	startRecord := (page-1)*perPage + 1
	endRecord := page * perPage
	if endRecord > int(total) {
		endRecord = int(total)
	}

	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// Renderizar vista
	c.HTML(http.StatusOK, "systems_list.html", gin.H{
		"title":            "Listado de Sistemas",
		"systems":          systems,
		"page":             page,
		"perPage":          perPage,
		"totalPages":       totalPages,
		"totalSystems":     total,
		"nameQuery":        nameQuery,
		"descriptionQuery": descQuery,
		"startRecord":      startRecord,
		"endRecord":        endRecord,
		"globals":          globals,
		"session":          sessionData.(middleware.SessionData),
		"navLink":          "systems",
	})
}
