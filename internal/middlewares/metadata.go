// internal/middlewares/file_validation.go
package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FileValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obtener configuraciones del .env
		maxSizeMB, _ := strconv.Atoi(os.Getenv("MAX_FILE_SIZE_MB"))                   // Ej: "5" para 5MB
		allowedExtensions := strings.Split(os.Getenv("ALLOWED_FILE_EXTENSIONS"), ",") // Ej: "jpg,png,pdf"
		// 2. Validar que las configuraciones existen
		if maxSizeMB == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Middleware Error, File size limit not configured", "message": "No se ha configurado el tamaño máximo del archivo.",
			})
			return
		}

		if len(allowedExtensions) == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Allowed extensions not configured", "message": "No se ha configurado las extensiones permitidas.",
			})
			return
		}

		// 3. Obtener el archivo del form-data
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "File is required", "message": "No se ha enviado archivo adjunto",
			})
			return
		}
		defer file.Close()

		// 4. Validar tamaño del archivo
		maxSizeBytes := int64(maxSizeMB) << 20 // Convertir MB a bytes
		if header.Size > maxSizeBytes {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("File too large. Max size is %dMB", maxSizeMB), "message": fmt.Sprintf("Archivo muy grande. Máximo permitido %dMB", maxSizeMB),
			})
			return
		}

		// 5. Validar extensión del archivo
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(header.Filename), "."))
		validExtension := false
		for _, allowedExt := range allowedExtensions {
			if ext == strings.ToLower(strings.TrimSpace(allowedExt)) {
				validExtension = true
				break
			}
		}

		if !validExtension {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Sólo se permiten las siguientes extensiones %s", strings.Join(allowedExtensions, ", ")), "error": fmt.Sprintf("Invalid file extension. Allowed: %s", strings.Join(allowedExtensions, ", ")),
			})
			return
		}

		// 6. Almacenar información del archivo en el contexto para handlers posteriores
		c.Set("fileInfo", gin.H{
			"filename": header.Filename,
			"size":     header.Size,
			"ext":      ext,
		})

		c.Next()
	}
}
