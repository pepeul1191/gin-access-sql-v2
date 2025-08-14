// main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"accessv2/internal/config"
	"accessv2/internal/middlewares"
	"accessv2/internal/routes"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Advertencia: No se encontró archivo .env o no se pudo cargar")
		// Puede continuar, usará valores por defecto
	}

	// Ahora puedes acceder a las variables de entorno
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000" // valor por defecto
	}

	// Asegurarse de que config cargue las variables (si usas config.Init())
	config.Init() // Opcional: si quieres inicializar config en un paquete

	// Crear router de Gin
	r := gin.Default()

	// 1. PRIMERO registrar el middleware CORS
	r.Use(middlewares.CORSMiddleware())

	// 2. LUEGO configurar rutas
	routes.Setup(r) // Asegúrate que Setup acepte *gin.Engine como parámetro

	// 3. Configurar archivos estáticos
	r.Static("/public", "./public")

	// Iniciar servidor
	fmt.Printf("🚀 Servidor escuchando en http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}
