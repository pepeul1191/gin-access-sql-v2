package systems

import (
	"accessv2/internal/handlers/permissions"
	"accessv2/internal/handlers/roles"
	"accessv2/internal/handlers/users"

	"github.com/gin-gonic/gin"
)

// He añadido 'userHandler' a los parámetros de la función.
func RegisterSystemsRoutes(r *gin.Engine, handler *SystemHandler, roleHandler *roles.RoleHandler, permissionHandler *permissions.PermissionHandler, userHandler *users.UserHandler) {

	// Main systems group
	systemsGroup := r.Group("/systems")
	{
		// Routes for listing and creating systems
		systemsGroup.GET("/", handler.ListSystems)
		systemsGroup.POST("/create", handler.CreateSystemHandler)
		systemsGroup.GET("/create", handler.CreateSystemHandler)

		// Group for a specific system identified by ':id'
		systemByIDGroup := systemsGroup.Group("/:id")
		{
			// Routes for editing and deleting a specific system
			systemByIDGroup.POST("/edit", handler.EditSystemHandler)
			systemByIDGroup.GET("/edit", handler.EditSystemHandler)
			systemByIDGroup.GET("/delete", handler.DeleteSystemHandler)

			// Routes for roles, now nested correctly under the specific system group
			systemByIDGroup.POST("/roles", roleHandler.CreateRoleHandler)
			systemByIDGroup.GET("/roles", roleHandler.CreateRoleHandler)
			systemByIDGroup.POST("/roles/:role_id/edit", roleHandler.EditRoleHandler)
			systemByIDGroup.GET("/roles/:role_id/edit", roleHandler.EditRoleHandler)
			systemByIDGroup.GET("/roles/:role_id/delete", roleHandler.DeleteRoleHandler)

			// permissions
			// Routes for roles, now nested correctly under the specific system group
			systemByIDGroup.GET("/roles/:role_id/permissions", handler.EditSystemHandler)
			systemByIDGroup.POST("/roles/:role_id/permissions/create", permissionHandler.CreatePermissionHandler)
			systemByIDGroup.GET("/roles/:role_id/permissions/create", permissionHandler.CreatePermissionHandler)
			systemByIDGroup.POST("/roles/:role_id/permissions/:permission_id/edit", permissionHandler.EditPermissionHandler)
			systemByIDGroup.GET("/roles/:role_id/permissions/:permission_id/edit", permissionHandler.EditPermissionHandler)
			systemByIDGroup.GET("/roles/:role_id/permissions/:permission_id/delete", permissionHandler.DeletePermissionHandler)

			//users
			systemByIDGroup.GET("/users", handler.ListSystemUsersHandler)
			systemByIDGroup.POST("/users", handler.SaveSystemUsersHandler)

			//users roles/permissions
			systemByIDGroup.GET("/users/:user_id", userHandler.GetUserRolesAndPermissions)
			systemByIDGroup.POST("/users/:user_id/permissions", userHandler.AssociatePermissionsHandler)
		}
	}
}
