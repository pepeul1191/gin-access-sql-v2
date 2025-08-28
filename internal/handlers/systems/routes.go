package systems

import (
	"accessv2/internal/handlers/roles"

	"github.com/gin-gonic/gin"
)

func RegisterSystemsRoutes(r *gin.Engine, handler *SystemHandler, roleHandler *roles.RoleHandler) {

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
		}
	}
}
