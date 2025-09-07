package main

import (
	"accessv2/config"
	"log"

	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	// 1. Configuración inicial
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}

	// 2. Inicialización de la base de datos
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// 4. Configuración de sesiones
	store := cookie.NewStore([]byte(config.GetEnv("SESSION_SECRET", "default-secret-32-bytes-long!")))

	// 5. Configuración del router
	router := config.SetupRouter(db, store)

	// 6. Configuración de vistas y estáticos
	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "./static") // Mejor nombre que 'public'

	// 7. Inicio del servidor
	serverPort := config.GetEnv("SERVER_PORT", ":8085")
	log.Printf("Server starting on %s", serverPort)
	if err := router.Run(serverPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
