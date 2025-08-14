// internal/handlers/common.go
package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	fmt.Println("Hola desde Gin!")
	c.String(200, "hol")
}
