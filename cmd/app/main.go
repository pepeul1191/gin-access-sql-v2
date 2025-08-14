// cmd/app/main.go
package main

import (
	"accessv2/config"

	"log"

	"github.com/gorilla/sessions"
)

func main() {
	// Inicializar base de datos
	err := config.InitDB()
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	// Configuración del session store
	store := sessions.NewCookieStore([]byte("tu-clave-secreta"))
	// Configuración de rutas
	router := config.SetupRouter(store)
	// Vistas
	router.LoadHTMLGlob("templates/**/*")
	// Archivos estáticos
	router.Static("/public", "./public")
	// Inicia servidor
	log.Println("Servidor iniciado en :8080")
	// Iniciar servidor
	router.Run(":8080")
}
