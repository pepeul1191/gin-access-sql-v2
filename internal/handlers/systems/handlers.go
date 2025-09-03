package systems

import (
	"accessv2/internal/domain"
	"accessv2/internal/forms"
	"accessv2/internal/services"
	"accessv2/pkg/middleware"
	"accessv2/pkg/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemHandler struct {
	service           *services.SystemService
	roleService       *services.RoleService
	permissionService *services.PermissionService
}

func NewSystemHandler(service *services.SystemService, roleService *services.RoleService, permissionService *services.PermissionService) *SystemHandler {
	return &SystemHandler{service: service, roleService: roleService, permissionService: permissionService}
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
	styles := []string{}
	scripts := []string{}

	// mensajes por URL, si lo hubiere
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}

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
		"styles":           styles,  // Pasar array de estilos
		"scripts":          scripts, // Pasar array de scripts
		"message":          message,
	})
}

func (h *SystemHandler) CreateSystemHandler(c *gin.Context) {
	// Obtener token CSRF una sola vez
	csrfToken, _ := c.Get("csrf_token")
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		var input forms.SystemCreateInput
		// Parsear formulario
		if err := c.ShouldBind(&input); err != nil {
			message := utils.Message{
				Content: err.Error(),
				Type:    "danger",
			}
			c.HTML(http.StatusBadRequest, "systems/create", gin.H{
				"title":   "Error al crear sistema",
				"error":   err.Error(),
				"csrf":    csrfToken,
				"values":  c.Request.PostForm,
				"globals": globals,
				"message": message,
				"session": sessionData.(middleware.SessionData),
				"navLink": "systems",
			})
			return
		}

		// Crear sistema a través del servicio
		system, err := h.service.CreateSystem(&input)
		if err != nil {
			message := utils.Message{
				Content: err.Error(),
				Type:    "danger",
			}
			c.HTML(http.StatusBadRequest, "systems/create", gin.H{
				"title":     "Error al crear sistema",
				"error":     err.Error(),
				"csrfToken": csrfToken,
				"values":    c.Request.PostForm,
				"globals":   globals,
				"message":   message,
				"session":   sessionData.(middleware.SessionData),
				"navLink":   "systems",
			})
			return
		}

		// Redirigir a editar sistema
		message := "Sistema creado exitosamente"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=success", system.ID, message))
		return
	}
	// Manejar método GET (muestra el formulario)
	c.HTML(http.StatusOK, "systems/create", gin.H{
		"title":     "Crear Nuevo Sistema",
		"globals":   globals,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "systems",
		"csrfToken": csrfToken,
	})
}

func (h *SystemHandler) EditSystemHandler(c *gin.Context) {
	// Obtener parámetros
	systemIdStr := c.Param("id")

	// Convertir el ID del sistema
	systemID, err := strconv.ParseUint(systemIdStr, 10, 32)
	if err != nil {
		message := "ID de sistema inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		h.handleEditSystemPost(c, systemID)
		return
	}

	// Verificar si el path contiene "permissions"
	if strings.Contains(c.Request.URL.Path, "permissions") {
		roleIdStr := c.Param("role_id")

		// Convertir el ID del sistema
		roleID, err := strconv.ParseUint(roleIdStr, 10, 32)
		if err != nil {
			message := "ID de rol inválido"
			c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
			return
		}

		// Lógica específica para permisos
		h.handleSystemRolesPermissions(c, systemID, roleID)
		return
	}

	// Manejar método GET (muestra el formulario)
	h.handleEditSystemGet(c, systemID)
}

func (h *SystemHandler) handleEditSystemPost(c *gin.Context, systemID uint64) {
	// Obtener datos del formulario
	form := forms.SystemEditInput{}

	if err := c.ShouldBind(&form); err != nil {
		message := "Datos del formulario inválidos"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=danger", systemID, url.QueryEscape(message)))
		return
	}

	// Validaciones adicionales
	if strings.TrimSpace(form.Name) == "" {
		message := "El nombre del sistema es requerido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=danger", systemID, url.QueryEscape(message)))
		return
	}

	// Obtener el sistema actual
	var system domain.System
	if err := h.service.FetchSystem(systemID, &system); err != nil {
		message := "Sistema no encontrado"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Actualizar datos
	system.Name = form.Name
	system.Description = form.Description
	system.Repository = form.Repository
	system.Updated = time.Now()

	// Guardar cambios
	if err := h.service.UpdateSystem(&system); err != nil {
		message := "Error al actualizar el sistema"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=danger", systemID, url.QueryEscape(message)))
		return
	}

	// Éxito - redireccionar con mensaje
	message := "Sistema actualizado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=success", systemID, url.QueryEscape(message)))
}

func (h *SystemHandler) handleEditSystemGet(c *gin.Context, systemID uint64) {
	// Obtener el sistema de la base de datos
	var system domain.System

	if err := h.service.FetchSystem(systemID, &system); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Sistema no encontrado"
		} else {
			message = "Error al cargar el sistema"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Obtener parámetros de paginación y búsqueda
	pageRole, _ := strconv.Atoi(c.DefaultQuery("page_roles", "1"))
	perPageRole, _ := strconv.Atoi(c.DefaultQuery("per_page_roles", "10"))

	var roles []domain.Role

	roles, totalRoles, err := h.roleService.GetPaginatedSystemRoles(pageRole, perPageRole, int(systemID))
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape("Error al buscar los roles del sistema")))
		return
	}

	// Calcular total de páginas
	totalPagesRoles := int(totalRoles) / perPageRole
	if int(totalRoles)%perPageRole > 0 {
		totalPagesRoles++
	}
	// Calcular registros mostrados
	startRecordRoles := (pageRole-1)*perPageRole + 1
	endRecordRoles := pageRole * perPageRole
	if endRecordRoles > int(totalRoles) {
		endRecordRoles = int(totalRoles)
	}

	// Obtener token CSRF
	csrfToken, _ := c.Get("csrf_token")
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// mensajes por URL, si lo hubiere
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}

	c.HTML(http.StatusOK, "systems/edit", gin.H{
		"title":            "Editar Sistema - " + system.Name,
		"csrfToken":        csrfToken,
		"globals":          globals,
		"system":           system,
		"session":          sessionData.(middleware.SessionData),
		"navLink":          "systems",
		"message":          message,
		"systemID":         systemID,
		"roles":            roles,
		"pageRole":         pageRole,
		"perPageRole":      perPageRole,
		"totalPagesRoles":  totalPagesRoles,
		"startRecordRoles": startRecordRoles,
		"endRecordRoles":   endRecordRoles,
		"totalRoles":       totalRoles,
		"styles":           []string{},
		"scripts":          []string{},
	})
}

func (h *SystemHandler) DeleteSystemHandler(c *gin.Context) {
	// Obtener parámetros
	systemIdStr := c.Param("id")

	// Convertir el ID del sistema
	systemID, err := strconv.ParseUint(systemIdStr, 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape("ID de sistema inválido")))
		return
	}

	// Verificar si el sistema existe primero
	var system domain.System
	if err := h.service.FetchSystem(systemID, &system); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Sistema no encontrado"
		} else {
			message = "Error al verificar el sistema"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Eliminar el sistema
	if err := h.service.DeleteSystem(systemID); err != nil {
		message := "Error al eliminar el sistema"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Éxito
	message := "Sistema eliminado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=success", url.QueryEscape(message)))
}

func (h *SystemHandler) handleSystemRolesPermissions(c *gin.Context, systemID uint64, roleID uint64) {
	// Obtener el sistema de la base de datos
	var system domain.System

	if err := h.service.FetchSystem(systemID, &system); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Sistema no encontrado"
		} else {
			message = "Error al cargar el sistema"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Obtener parámetros de paginación y búsqueda
	pageRole, _ := strconv.Atoi(c.DefaultQuery("page_roles", "1"))
	perPageRole, _ := strconv.Atoi(c.DefaultQuery("per_page_roles", "10"))

	var roles []domain.Role

	roles, totalRoles, err := h.roleService.GetPaginatedSystemRoles(pageRole, perPageRole, int(systemID))
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape("Error al buscar los roles del sistema")))
		return
	}

	// Calcular total de páginas
	totalPagesRoles := int(totalRoles) / perPageRole
	if int(totalRoles)%perPageRole > 0 {
		totalPagesRoles++
	}
	// Calcular registros mostrados
	startRecordRoles := (pageRole-1)*perPageRole + 1
	endRecordRoles := pageRole * perPageRole
	if endRecordRoles > int(totalRoles) {
		endRecordRoles = int(totalRoles)
	}

	// buscar rol por id
	var role domain.Role
	if err := h.roleService.FetchRole(roleID, &role); err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape("Error al buscar los permisos del rol")))
		return
	}

	// Obtener parámetros de paginación y búsqueda del permiso
	pagePermission, _ := strconv.Atoi(c.DefaultQuery("page_permissions", "1"))
	perPagePermission, _ := strconv.Atoi(c.DefaultQuery("per_page_permissions", "10"))

	var permissions []domain.Permission

	permissions, totalPermissions, err := h.permissionService.GetPaginatedRolePermissions(pagePermission, perPagePermission, int(roleID))
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d?message=%s&type=danger", url.QueryEscape("Error al buscar los roles del sistema")))
		return
	}

	// Calcular total de páginas
	totalPagesPermissions := int(totalPermissions) / perPagePermission
	if int(totalRoles)%perPageRole > 0 {
		totalPagesPermissions++
	}
	// Calcular registros mostrados
	startRecordPermissions := (pagePermission-1)*perPagePermission + 1
	endRecordPermissions := pagePermission * perPagePermission
	if endRecordPermissions > int(totalPermissions) {
		endRecordPermissions = int(totalPermissions)
	}

	// Obtener token CSRF
	csrfToken, _ := c.Get("csrf_token")
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// mensajes por URL, si lo hubiere
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}

	c.HTML(http.StatusOK, "systems/permissions", gin.H{
		"title":     "Editar Sistema - " + system.Name,
		"csrfToken": csrfToken,
		"globals":   globals,
		"system":    system,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "systems",
		"message":   message,
		"systemID":  systemID,
		// roles
		"roles":            roles,
		"role":             role,
		"pageRole":         pageRole,
		"perPageRole":      perPageRole,
		"totalPagesRoles":  totalPagesRoles,
		"startRecordRoles": startRecordRoles,
		"endRecordRoles":   endRecordRoles,
		"totalRoles":       totalRoles,
		// permissions
		"roleID":                 roleID,
		"permissions":            permissions,
		"pagePermission":         pagePermission,
		"perPagePermission":      perPagePermission,
		"totalPagesPermissions":  totalPagesPermissions,
		"startRecordPermissions": startRecordPermissions,
		"endRecordPermissions":   endRecordPermissions,
		"totalPermissions":       totalPermissions,
		"styles":                 []string{},
		"scripts":                []string{},
	})
}

func (h *SystemHandler) ListSystemUsersHandler(c *gin.Context) {
	// Obtener parámetros
	systemIdStr := c.Param("id")

	// Convertir el ID del sistema
	systemID, err := strconv.ParseUint(systemIdStr, 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape("ID de sistema inválido")))
		return
	}

	// Obtener parámetros de paginación y búsqueda
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	usernameQuery := strings.TrimSpace(c.Query("username"))
	emailQuery := strings.TrimSpace(c.Query("email"))
	statusQuery := strings.TrimSpace(c.DefaultQuery("association_status", "2"))

	// Validar parámetros
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Obtener usuarios paginados
	users, total, err := h.service.GetPaginatedSystemUsers(page, perPage, usernameQuery, emailQuery, statusQuery, systemID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "Error al obtener los usuarios",
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
	styles := []string{}
	scripts := []string{}

	// mensajes por URL, si lo hubiere
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}

	// Renderizar vista
	c.HTML(http.StatusOK, "systems/users", gin.H{
		"title":         "Usuarios del Sistemas",
		"users":         users,
		"page":          page,
		"perPage":       perPage,
		"totalPages":    totalPages,
		"totalUsers":    total,
		"usernameQuery": usernameQuery,
		"emailQuery":    emailQuery,
		"statusQuery":   statusQuery,
		"systemID":      systemID,
		"startRecord":   startRecord,
		"endRecord":     endRecord,
		"globals":       globals,
		"session":       sessionData.(middleware.SessionData),
		"navLink":       "systems",
		"styles":        styles,  // Pasar array de estilos
		"scripts":       scripts, // Pasar array de scripts
		"message":       message,
	})
}
