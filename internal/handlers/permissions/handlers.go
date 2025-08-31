package permissions

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

type PermissionHandler struct {
	service *services.PermissionService
}

func NewPermissionHandler(service *services.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	// Obtener parámetros de paginación y búsqueda
	// Obtener parámetros
	roleIdStr := c.Param("id")
	// Convertir el ID del usuario
	roleID, err := strconv.ParseInt(roleIdStr, 10, 32)
	if err != nil {
		message := "ID de rol inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Obtener sistemas paginados
	permissions, err := h.service.GetAllByRoleID(int(roleID))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "Error al obtener los usuarios",
		})
		return
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
		"title":       "Listado de Usuarios",
		"permissions": permissions,
		"roleID":      roleID,
		"globals":     globals,
		"session":     sessionData.(middleware.SessionData),
		"navLink":     "systems",
		"styles":      styles,  // Pasar array de estilos
		"scripts":     scripts, // Pasar array de scripts
		"message":     message,
	})
}

func (h *PermissionHandler) CreatePermissionHandler(c *gin.Context) {
	// Obtener token CSRF una sola vez
	csrfToken, _ := c.Get("csrf_token")
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// Obtener parámetros
	systemIdStr := c.Param("id")
	roleIdStr := c.Param("role_id")
	// Convertir el ID del usuario
	systemID, err := strconv.ParseInt(systemIdStr, 10, 32)
	if err != nil {
		message := "ID de sistema inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	roleID, err := strconv.ParseInt(roleIdStr, 10, 32)
	if err != nil {
		message := "ID de rol inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions/create?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		var input forms.PermissionCreateInput
		// Parsear formulario
		if err := c.ShouldBind(&input); err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions/create?message=%s&type=danger", systemID, roleID, err.Error()))
			return
		}

		// Crear usuario a través del servicio
		permission, err := h.service.CreatePermission(&input, int(roleID))
		if err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions/create?message=%s&type=danger", systemID, roleID, err.Error()))
			return
		}

		// Redirigir a editar usuario
		message := "Permis creado exitosamente"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions/%d/edit?message=%s&type=success", systemID, roleID, permission.ID, message))
		return
	}
	// Manejar método GET (muestra el formulario)
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}
	c.HTML(http.StatusOK, "permissions/create", gin.H{
		"title":     "Crear Nuevo Permiso",
		"systemID":  systemID,
		"roleID":    roleID,
		"globals":   globals,
		"message":   message,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "systems",
		"csrfToken": csrfToken,
	})
}

func (h *PermissionHandler) EditPermissionHandler(c *gin.Context) {
	// Obtener parámetros
	systemIdStr := c.Param("id")
	roleIdStr := c.Param("role_id")
	permissionIdStr := c.Param("permission_id")

	// Convertir el ID del sistema
	systemID, err := strconv.ParseUint(systemIdStr, 10, 32)
	if err != nil {
		message := "ID de sistema no inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Convertir el ID del rol
	roleID, err := strconv.ParseUint(roleIdStr, 10, 32)
	if err != nil {
		message := "ID de rol no inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=danger", systemID, url.QueryEscape(message)))
		return
	}

	// Convertir el ID del permiso
	permissionID, err := strconv.ParseUint(permissionIdStr, 10, 32)
	if err != nil {
		message := "ID de permiso no inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permission/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		h.handleEditPermissionPost(c, systemID, roleID, permissionID)
		return
	}

	// Manejar método GET (muestra el formulario)
	h.handleEditPermissionGet(c, systemID, roleID, permissionID)
}

func (h *PermissionHandler) handleEditPermissionPost(c *gin.Context, systemID uint64, roleID uint64, permissionID uint64) {
	// Obtener datos del formulario
	form := forms.PermissionEditInput{}

	if err := c.ShouldBind(&form); err != nil {
		message := "Datos del formulario inválidos"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions/%d/edit?message=%s&type=danger", systemID, roleID, permissionID, url.QueryEscape(message)))
		return
	}

	// Validaciones adicionales
	if strings.TrimSpace(form.Name) == "" {
		message := "El nombre del permiso es requerido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions/%d/edit?message=%s&type=danger", systemID, roleID, permissionID, url.QueryEscape(message)))
		return
	}

	// Obtener el usuario actual
	var permission domain.Permission
	if err := h.service.FetchPermission(permissionID, &permission); err != nil {
		message := "Permiso no encontrado"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Actualizar datos
	permission.Name = form.Name
	permission.Updated = time.Now()

	// Guardar cambios
	if err := h.service.UpdatePermssion(&permission); err != nil {
		message := err.Error()
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Éxito - redireccionar con mensaje
	message := "Permiso actualizado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d?message=%s&type=success", systemID, roleID, url.QueryEscape(message)))
}

func (h *PermissionHandler) handleEditPermissionGet(c *gin.Context, systemID uint64, roleID uint64, permissionID uint64) {
	// Obtener el sistema de la base de datos
	var permission domain.Permission

	if err := h.service.FetchPermission(permissionID, &permission); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Permiso no encontrado"
		} else {
			message = "Error al cargar el permiso"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
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

	c.HTML(http.StatusOK, "permissions/edit", gin.H{
		"title":      "Editar Sistema",
		"csrfToken":  csrfToken,
		"globals":    globals,
		"permission": permission,
		"systemID":   systemID,
		"roleID":     roleID,
		"session":    sessionData.(middleware.SessionData),
		"navLink":    "systems",
		"message":    message,
		"styles":     []string{},
		"scripts":    []string{},
	})
}

func (h *PermissionHandler) DeletePermissionHandler(c *gin.Context) {
	// Obtener parámetros
	systemIdStr := c.Param("id")
	roleIdStr := c.Param("role_id")
	permissionIdStr := c.Param("permission_id")

	// Convertir el ID del sistema
	systemID, err := strconv.ParseUint(systemIdStr, 10, 32)
	if err != nil {
		message := "ID de sistema no inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Convertir el ID del rol
	roleID, err := strconv.ParseUint(roleIdStr, 10, 32)
	if err != nil {
		message := "ID de rol no inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Convertir el ID del rol
	permissionID, err := strconv.ParseUint(permissionIdStr, 10, 32)
	if err != nil {
		message := "ID de permiso no inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Verificar si el rol existe primero
	var permission domain.Permission
	if err := h.service.FetchPermission(permissionID, &permission); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Permisso no encontrado"
		} else {
			message = "Error al verificar el permiso"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Eliminar el permiso
	if err := h.service.DeletePermission(permissionID); err != nil {
		message := "Error al eliminar el rol"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Éxito
	message := "Permiso eliminado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/permissions?message=%s&type=success", systemID, roleID, url.QueryEscape(message)))
}
