// cmd/app/main.go
package main

import (
	"accessv2/config"
	"log"
)

func main() {
	// Configuración de rutas
	router := config.SetupRouter()
	// Vistas
	router.LoadHTMLGlob("templates/*.html")
	// Inicia servidor
	log.Println("Servidor iniciado en :8080")
	// Iniciar servidor
	router.Run(":8080")
}
