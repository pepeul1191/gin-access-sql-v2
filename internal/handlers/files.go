// internal/handlers/files.go
package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const UploadsDir = "uploads"
const PublicDir = "public"

func UploadFile(c *gin.Context) {
	foldeName := c.Param("folder_name")

	file, header, err := c.Request.FormFile("file")
	if file == nil || err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "No se proporcionó ningún archivo"})
		return
	}
	defer file.Close()

	dir := filepath.Join(UploadsDir, foldeName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "No se pudo crear el directorio"})
		return
	}

	ext := filepath.Ext(header.Filename)
	randomName := fmt.Sprintf("%x%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, randomName)

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error al crear el archivo"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error al guardar el archivo"})
		return
	}

	// Obtener MIME type
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Si no se reconoce por extensión, intentar detectarlo
		mimeType = http.DetectContentType([]byte(ext))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"filename":          randomName,
		"path":              fmt.Sprintf("/uploads/%s/%s", foldeName, randomName),
		"original_filename": header.Filename,
		"size":              header.Size,
		"mime_type":         mimeType,
	})
}

func UploadFileToPublic(c *gin.Context) {
	foldeName := c.Param("folder_name")

	file, header, err := c.Request.FormFile("file")
	if file == nil || err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "No se proporcionó ningún archivo"})
		return
	}
	defer file.Close()

	dir := filepath.Join(PublicDir, UploadsDir, foldeName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "No se pudo crear el directorio"})
		return
	}

	ext := filepath.Ext(header.Filename)
	randomName := fmt.Sprintf("%x%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(dir, randomName)

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error al crear el archivo"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error al guardar el archivo"})
		return
	}

	// Obtener MIME type
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// Si no se reconoce por extensión, intentar detectarlo
		mimeType = http.DetectContentType([]byte(ext))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            "success",
		"filename":          randomName,
		"path":              fmt.Sprintf("/uploads/%s/%s", foldeName, randomName),
		"original_filename": header.Filename,
		"size":              header.Size,
		"mime_type":         mimeType,
	})
}

func DownloadFile(c *gin.Context) {
	foldeName := c.Param("folder_name")
	fileName := c.Param("file_name")

	filePath := filepath.Join(UploadsDir, foldeName, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "message": "File not found or not accessible"})
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Header("Content-Type", contentType)

	c.File(filePath)
}
