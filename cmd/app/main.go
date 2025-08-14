// cmd/app/main.go
package main

import (
	"accessv2/config"

	"log"
)

func main() {
	// Inicializar base de datos
	err := config.InitDB()
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	// Configuración de rutas
	router := config.SetupRouter()
	// Vistas
	router.LoadHTMLGlob("templates/**/*")
	// Inicia servidor
	log.Println("Servidor iniciado en :8080")
	// Iniciar servidor
	router.Run(":8080")
}
