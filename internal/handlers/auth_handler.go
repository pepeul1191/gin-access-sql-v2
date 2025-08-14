// internal/handlers/auth_handler.go
package handlers

import (
	"accessv2/internal/domain"
	"accessv2/internal/forms"
	"accessv2/internal/services"
	"html/template"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
	tmpl        *template.Template
}

func NewAuthHandler(authService *services.AuthService, tmpl *template.Template) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		tmpl:        tmpl,
	}
}

// ShowRegister muestra el formulario de registro (GET)
func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Form  *forms.RegisterForm
		Error string
	}{
		Form:  &forms.RegisterForm{},
		Error: "",
	}
	h.tmpl.ExecuteTemplate(w, "register.html", data)
}

// Register procesa el formulario de registro (POST)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parsear el formulario
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Crear el formulario con los datos recibidos
	form := &forms.RegisterForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	// Validar el formulario
	if form.Username == "" || form.Password == "" {
		h.tmpl.ExecuteTemplate(w, "register.html", struct {
			Form  *forms.RegisterForm
			Error string
		}{
			Form:  form,
			Error: "Username y password son requeridos",
		})
		return
	}

	// Registrar al usuario
	err = h.authService.Register(&domain.User{
		Username: form.Username,
		Password: form.Password, // Deberías hashear la contraseña aquí
	})

	if err != nil {
		h.tmpl.ExecuteTemplate(w, "register.html", struct {
			Form  *forms.RegisterForm
			Error string
		}{
			Form:  form,
			Error: "Error al registrar: " + err.Error(),
		})
		return
	}

	// Redirigir al login después del registro exitoso
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
