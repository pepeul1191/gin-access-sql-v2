// internal/handlers/auth.go
package handlers

import (
	"fmt"
	"net/http"

	"accessv2/internal/utils"

	"github.com/gin-gonic/gin"
)

func SignIn(c *gin.Context) {
	fmt.Println("Authentication successful. Generating JWT.")
	token, err := utils.GenerateJWT()
	if err != nil {
		fmt.Println("Failed to encode JWT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
