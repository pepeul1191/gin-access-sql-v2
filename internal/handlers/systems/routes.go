package systems

import (
	"github.com/gin-gonic/gin"
)

func RegisterSystemsRoutes(r *gin.Engine, handler *SystemHandler) {

	r.GET("/systems", handler.ListSystems)
}
