// cmd/app/main.go
package main

import (
	"accessv2/config"

	"log"

	"github.com/gin-contrib/sessions/cookie" // Importar el store de cookies
)

func main() {
	// Inicializar base de datos
	err := config.InitDB()
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}

	// Configuración del session store usando gin-contrib/sessions
	store := cookie.NewStore([]byte("clave-secreta-de-32-bytes"))
	// Configuración de rutas
	router := config.SetupRouter(store)

	// Configuración de vistas y archivos estáticos
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/public", "./public")

	// Inicia servidor
	log.Println("Servidor iniciado en :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
