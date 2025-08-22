package systems

import (
	"github.com/gin-gonic/gin"
)

func RegisterSystemsRoutes(r *gin.Engine, handler *SystemHandler) {

	r.GET("/systems", handler.ListSystems)
	r.POST("/systems/create", handler.CreateSystemHandler)
	r.GET("/systems/create", handler.CreateSystemHandler)
	r.POST("/systems/:id/edit", handler.EditSystemHandler)
	r.GET("/systems/:id/edit", handler.EditSystemHandler)
	r.GET("/systems/:id/delete", handler.DeleteSystemHandler)
}
