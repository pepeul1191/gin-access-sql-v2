package roles

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

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(service *services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	// Obtener parámetros de paginación y búsqueda
	// Obtener parámetros
	systemIdStr := c.Param("id")
	// Convertir el ID del usuario
	systemID, err := strconv.ParseInt(systemIdStr, 10, 32)
	if err != nil {
		message := "ID de sistema inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Obtener sistemas paginados
	roles, err := h.service.GetAllBySystemID(int(systemID))
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
		"title":    "Listado de Usuarios",
		"roles":    roles,
		"systemID": systemID,
		"globals":  globals,
		"session":  sessionData.(middleware.SessionData),
		"navLink":  "systems",
		"styles":   styles,  // Pasar array de estilos
		"scripts":  scripts, // Pasar array de scripts
		"message":  message,
	})
}

func (h *RoleHandler) CreateRoleHandler(c *gin.Context) {
	// Obtener token CSRF una sola vez
	csrfToken, _ := c.Get("csrf_token")
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// Obtener parámetros
	systemIdStr := c.Param("id")
	// Convertir el ID del usuario
	systemID, err := strconv.ParseInt(systemIdStr, 10, 32)
	if err != nil {
		message := "ID de sistema inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems?message=%s&type=danger", url.QueryEscape(message)))
		return
	}
	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		var input forms.RoleCreateInput
		// Parsear formulario
		if err := c.ShouldBind(&input); err != nil {
			message := utils.Message{
				Content: err.Error(),
				Type:    "danger",
			}
			csrfToken, _ := c.Get("csrf_token")
			c.HTML(http.StatusBadRequest, fmt.Sprintf("roles/create"), gin.H{
				"title":     "Error al crear usuario",
				"error":     err.Error(),
				"csrfToken": csrfToken,
				"form":      c.Request.PostForm,
				"globals":   globals,
				"message":   message,
				"systemID":  systemID,
				"session":   sessionData.(middleware.SessionData),
				"navLink":   "users",
			})
			return
		}

		// Crear usuario a través del servicio
		role, err := h.service.CreateRole(&input, int(systemID))
		if err != nil {
			message := utils.Message{
				Content: err.Error(),
				Type:    "danger",
			}
			c.HTML(http.StatusBadRequest, fmt.Sprintf("systems/%d/edit/", systemID), gin.H{
				"title":     "Error al crear rol",
				"error":     err.Error(),
				"csrfToken": csrfToken,
				"form":      input,
				"globals":   globals,
				"message":   message,
				"session":   sessionData.(middleware.SessionData),
				"navLink":   "users",
			})
			return
		}

		// Redirigir a editar usuario
		message := "Rol creado exitosamente"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=success", systemID, role.ID, message))
		return
	}
	// Manejar método GET (muestra el formulario)
	c.HTML(http.StatusOK, "roles/create", gin.H{
		"title":     "Crear Nuevo Usuario",
		"systemID":  systemID,
		"globals":   globals,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "users",
		"csrfToken": csrfToken,
	})
}

func (h *RoleHandler) EditRoleHandler(c *gin.Context) {
	// Obtener parámetros
	systemIdStr := c.Param("id")
	roleIdStr := c.Param("role_id")

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

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		h.handleEditRolePost(c, systemID, roleID)
		return
	}

	// Manejar método GET (muestra el formulario)
	h.handleEditRoleGet(c, systemID, roleID)
}

func (h *RoleHandler) handleEditRolePost(c *gin.Context, systemID uint64, roleID uint64) {
	// Obtener datos del formulario
	form := forms.RoleEditInput{}

	if err := c.ShouldBind(&form); err != nil {
		message := "Datos del formulario inválidos"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Validaciones adicionales
	if strings.TrimSpace(form.Name) == "" {
		message := "El nombre del rol es requerido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Obtener el usuario actual
	var role domain.Role
	if err := h.service.FetchRole(roleID, &role); err != nil {
		message := "Rol no encontrado"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=danger", systemID, url.QueryEscape(message)))
		return
	}

	// Actualizar datos
	role.Name = form.Name
	role.Updated = time.Now()

	// Guardar cambios
	if err := h.service.UpdateRole(&role); err != nil {
		message := "Error al actualizar el rol"
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=danger", systemID, roleID, url.QueryEscape(message)))
		return
	}

	// Éxito - redireccionar con mensaje
	message := "Rol actualizado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/roles/%d/edit?message=%s&type=success", systemID, roleID, url.QueryEscape(message)))
}

func (h *RoleHandler) handleEditRoleGet(c *gin.Context, systemID uint64, roleID uint64) {
	// Obtener el sistema de la base de datos
	var role domain.Role

	if err := h.service.FetchRole(roleID, &role); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Rol no encontrado"
		} else {
			message = "Error al cargar el rol"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/edit?message=%s&type=danger", systemID, url.QueryEscape(message)))
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

	c.HTML(http.StatusOK, "roles/edit", gin.H{
		"title":     "Editar Sistema",
		"csrfToken": csrfToken,
		"globals":   globals,
		"role":      role,
		"systemID":  systemID,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "systems",
		"message":   message,
		"styles":    []string{},
		"scripts":   []string{},
	})
}

/*
func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	// Obtener parámetros
	userIDStr := c.Param("id")

	// Convertir el ID del usuario
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape("ID de usuario inválido")))
		return
	}

	// Verificar si el usuario existe primero
	var user domain.User
	if err := h.service.FetchUser(userID, &user); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Usuario no encontrado"
		} else {
			message = "Error al verificar el usuario"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Eliminar el usuario
	if err := h.service.DeleteUser(userID); err != nil {
		message := "Error al eliminar el usuario"
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Éxito
	message := "Ysuario eliminado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=success", url.QueryEscape(message)))
}
*/
