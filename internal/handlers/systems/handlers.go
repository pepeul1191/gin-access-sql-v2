package systems

import (
	"accessv2/internal/forms"
	"accessv2/internal/services"
	"accessv2/pkg/middleware"
	"fmt"
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
	c.HTML(http.StatusOK, "systems/list", gin.H{
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

func (h *SystemHandler) CreateSystemHandler(c *gin.Context) {
	// Obtener token CSRF una sola vez
	csrfToken := c.MustGet("csrf_token").(string)
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")
	fmt.Println("1 ++++++++++++++++++++")

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		var input forms.SystemCreateInput

		// Parsear formulario
		if err := c.ShouldBind(&input); err != nil {
			c.HTML(http.StatusBadRequest, "systems/create", gin.H{
				"title":   "Error al crear sistema",
				"error":   "Datos inválidos",
				"csrf":    csrfToken,
				"values":  c.Request.PostForm,
				"globals": globals,
				"session": sessionData.(middleware.SessionData),
				"navLink": "systems",
			})
			return
		}

		// Crear sistema a través del servicio
		system, err := h.service.CreateSystem(&input)
		if err != nil {
			c.HTML(http.StatusBadRequest, "systems/create", gin.H{
				"title":   "Error al crear sistema",
				"error":   err.Error(),
				"csrf":    csrfToken,
				"values":  c.Request.PostForm,
				"globals": globals,
				"session": sessionData.(middleware.SessionData),
				"navLink": "systems",
			})
			return
		}

		// Redirigir al listado con mensaje de éxito
		c.Redirect(http.StatusFound, "/systems?success=Sistema creado exitosamente: "+system.Name)
		return
	}
	fmt.Println("2 ++++++++++++++++++++")
	// Manejar método GET (muestra el formulario)
	c.HTML(http.StatusOK, "systems/create", gin.H{
		"title":   "Crear Nuevo Sistema",
		"globals": globals,
		"session": sessionData.(middleware.SessionData),
		"navLink": "systems",
		"csrf":    csrfToken,
	})
}
