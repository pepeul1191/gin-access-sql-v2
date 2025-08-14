// cmd/app/main.go
package main

import (
	"accessv2/config"
	"accessv2/internal/handlers"
	"accessv2/internal/repositories"
	"accessv2/internal/services"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Configuración
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Templates
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Inicialización de dependencias
	userRepo := repositories.NewUserRepository(cfg.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, tmpl)

	// Rutas
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.ShowRegister(w, r)
		case http.MethodPost:
			authHandler.Register(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Inicia servidor
	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
