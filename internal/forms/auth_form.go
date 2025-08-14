// internal/forms/auth_form.go
package forms

type RegisterForm struct {
	Username string `form:"username" binding:"required,min=3,max=20"`
	Password string `form:"password" binding:"required,min=8"`
	Email    string `form:"email" binding:"omitempty,email"` // Opcional pero debe ser email válido si existe
}

// Validate puede añadirse para lógica personalizada
func (f *RegisterForm) Validate() error {
	// Ejemplo: Verificar que el username no contenga espacios
	// if strings.Contains(f.Username, " ") {
	//     return errors.New("username no puede contener espacios")
	// }
	return nil
}
