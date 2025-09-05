package users

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

type UserHandler struct {
	service               *services.UserService
	userPermissionService *services.UserPermissionService
}

func NewUserHandler(service *services.UserService, userPermissionService *services.UserPermissionService) *UserHandler {
	return &UserHandler{service: service, userPermissionService: userPermissionService}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	// Obtener parámetros de paginación y búsqueda
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	usernameQuery := strings.TrimSpace(c.Query("username"))
	emailQuery := strings.TrimSpace(c.Query("email"))
	statusQuery := strings.TrimSpace(c.Query("status")) // Nuevo parámetro

	// Validar parámetros
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Obtener sistemas paginados
	users, total, err := h.service.GetPaginatedUsers(page, perPage, usernameQuery, emailQuery, statusQuery)
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
	c.HTML(http.StatusOK, "users/list", gin.H{
		"title":         "Listado de Usuarios",
		"users":         users,
		"page":          page,
		"perPage":       perPage,
		"totalPages":    totalPages,
		"totalUsers":    total,
		"usernameQuery": usernameQuery,
		"emailQuery":    emailQuery,
		"statusQuery":   statusQuery, // Nuevo parámetro
		"startRecord":   startRecord,
		"endRecord":     endRecord,
		"globals":       globals,
		"session":       sessionData.(middleware.SessionData),
		"navLink":       "users",
		"styles":        styles,  // Pasar array de estilos
		"scripts":       scripts, // Pasar array de scripts
		"message":       message,
	})
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	// Obtener token CSRF una sola vez
	csrfToken, _ := c.Get("csrf_token")
	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		var input forms.UserCreateInput
		// Parsear formulario
		if err := c.ShouldBind(&input); err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/users/create?message=%s&type=danger", err.Error()))
			return
		}

		// Crear usuario a través del servicio
		user, err := h.service.CreateUser(&input)
		if err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/users/create?message=%s&type=danger", err.Error()))
			return
		}

		// Redirigir a editar usuario
		message := "Usuario creado exitosamente"
		c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d/edit?message=%s&type=success", user.ID, message))
		return
	}
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}
	// Manejar método GET (muestra el formulario)
	c.HTML(http.StatusOK, "users/create", gin.H{
		"title":     "Crear Nuevo Usuario",
		"globals":   globals,
		"message":   message,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "users",
		"csrfToken": csrfToken,
	})
}

func (h *UserHandler) EditUserHandler(c *gin.Context) {
	// Obtener parámetros
	userIdStr := c.Param("id")

	// Convertir el ID del usuario
	userID, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		message := "ID de usuario inválido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Manejar método POST
	if c.Request.Method == http.MethodPost {
		h.handleEditUserPost(c, userID)
		return
	}

	// Manejar método GET (muestra el formulario)
	h.handleEditUserGet(c, userID)
}

func (h *UserHandler) handleEditUserPost(c *gin.Context, userID uint64) {
	// Obtener datos del formulario
	form := forms.UserEditInput{}

	if err := c.ShouldBind(&form); err != nil {
		message := "Datos del formulario inválidos"
		c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d/edit?message=%s&type=danger", userID, url.QueryEscape(message)))
		return
	}

	// Validaciones adicionales
	if strings.TrimSpace(form.Username) == "" {
		message := "El nombre del usuario es requerido"
		c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d/edit?message=%s&type=danger", userID, url.QueryEscape(message)))
		return
	}

	// Obtener el usuario actual
	var user domain.User
	if err := h.service.FetchUser(userID, &user); err != nil {
		message := "Usuario no encontrado"
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape(message)))
		return
	}

	// Actualizar datos
	user.Username = form.Username
	user.Email = form.Email
	if form.Password != "1234567890" {
		user.Password = form.Password
	}
	if form.Status == "active" {
		user.Activated = true
	} else {
		user.Activated = false
	}
	user.Updated = time.Now()

	// Guardar cambios
	if err := h.service.UpdateUser(&user); err != nil {
		message := err.Error()
		c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d/edit?message=%s&type=danger", userID, url.QueryEscape(message)))
		return
	}

	// Éxito - redireccionar con mensaje
	message := "Usuario actualizado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/users/%d/edit?message=%s&type=success", userID, url.QueryEscape(message)))
}

func (h *UserHandler) handleEditUserGet(c *gin.Context, userID uint64) {
	// Obtener el sistema de la base de datos
	var user domain.User

	if err := h.service.FetchUser(userID, &user); err != nil {
		message := ""
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "Usuario no encontrado"
		} else {
			message = "Error al cargar el usuario"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape(message)))
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

	// cambiar contraseña
	user.Password = "1234567890"

	c.HTML(http.StatusOK, "users/edit", gin.H{
		"title":     "Editar Sistema",
		"csrfToken": csrfToken,
		"globals":   globals,
		"user":      user,
		"session":   sessionData.(middleware.SessionData),
		"navLink":   "users",
		"message":   message,
		"styles":    []string{},
		"scripts":   []string{},
	})
}

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
	message := "Usuario eliminado exitosamente"
	c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=success", url.QueryEscape(message)))
}

func (h *UserHandler) GetUserRolesAndPermissions(c *gin.Context) {
	// Obtener parámetros
	userIDStr := c.Param("user_id")
	systemIDStr := c.Param("id")

	// Convertir el ID del sistema
	systemID, err := strconv.ParseUint(systemIDStr, 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape("ID de sistema inválido")))
		return
	}

	// Convertir el ID del usuario
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/users?message=%s&type=danger", url.QueryEscape("ID de usuario inválido")))
		return
	}

	// Obtener las relaciones de roles y permisos
	permissions, err := h.userPermissionService.GetUserRolesAndPermissions(systemID, userID)
	if err != nil {
		// Aquí puedes manejar el error de forma más específica si es necesario,
		// por ejemplo, c.Status(http.StatusInternalServerError).
		return
	}

	globals, _ := c.Get("globals")
	sessionData, _ := c.Get("sessionData")
	styles := []string{}
	scripts := []string{}
	csrfToken, _ := c.Get("csrf_token")

	// mensajes por URL, si lo hubiere
	message := utils.Message{
		Content: c.Query("message"),
		Type:    c.Query("type"),
	}

	// Renderizar vista
	c.HTML(http.StatusOK, "users/roles-permissions", gin.H{
		"title":       "Permisos de los Roles del Usuario",
		"systemID":    systemID,
		"userID":      userID,
		"permissions": permissions,
		"csrfToken":   csrfToken,
		"globals":     globals,
		"session":     sessionData.(middleware.SessionData),
		"navLink":     "systems",
		"styles":      styles,  // Pasar array de estilos
		"scripts":     scripts, // Pasar array de scripts
		"message":     message,
	})
}

func (h *UserHandler) AssociatePermissionsHandler(c *gin.Context) {
	// Recuperar los parámetros de la URL (systemID y userID)
	systemIDStr := c.Param("id")
	userIDStr := c.Param("user_id")

	// Convertir los IDs a uint64
	systemID, err := strconv.ParseUint(systemIDStr, 10, 64)
	if err != nil {
		// Redirigir con un mensaje de error
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/?message=%s&type=danger", url.QueryEscape("ID de sistema inválido")))
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		// Redirigir con un mensaje de error
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/users?message=%s&type=danger", systemID, url.QueryEscape("ID de usuario inválido")))
		return
	}

	// Obtener los permisos seleccionados del formulario (parámetros de tipo checkbox)
	permissions := c.PostFormMap("permissions")

	// Convertir los valores de los permisos seleccionados a enteros (IDs)
	var permissionIDs []uint64
	for permIDStr := range permissions {
		permID, err := strconv.ParseUint(permIDStr, 10, 64)
		if err != nil {
			// Redirigir con un mensaje de error
			c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/users/%d?message=%s&type=danger", systemID, userID, url.QueryEscape("Error al procesar los permisos")))
			return
		}
		permissionIDs = append(permissionIDs, permID)
	}

	// Llamar a un servicio o repositorio para asociar los permisos al usuario
	err = h.userPermissionService.AssociatePermissions(uint(systemID), uint(userID), permissionIDs)
	if err != nil {
		// Redirigir con un mensaje de error
		c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/users/%d?message=%s&type=danger", systemID, userID, url.QueryEscape("Error al asociar los permisos")))
		return
	}

	// Redirigir al usuario de vuelta con un mensaje de éxito
	c.Redirect(http.StatusFound, fmt.Sprintf("/systems/%d/users/%d?message=%s&type=success", systemID, userID, url.QueryEscape("Permisos actualizados con éxito")))
}
