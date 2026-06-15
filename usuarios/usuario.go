// Package usuarios gestiona el registro, autenticación y perfiles
// de los usuarios del sistema de streaming.
package usuarios

import (
	"errors"
	"fmt"

	"streaming-system/suscripcion"
)

// Usuario representa a una persona registrada en el sistema.
//
// ENCAPSULACIÓN: todos los campos son privados (minúscula). En
// especial, "password" nunca puede leerse ni modificarse directamente
// desde fuera del paquete — solo a través de métodos como Login()
// o CambiarPassword().
type Usuario struct {
	id       int
	nombre   string
	email    string
	password string
	plan     suscripcion.Plan
	activo   bool
}

// contador interno para generar IDs automáticamente.
var contadorID int = 1

// NuevoUsuario registra un nuevo usuario en el sistema.
//
// MANEJO DE ERRORES: valida nombre, email y longitud mínima de la
// contraseña antes de crear el usuario. Si algo falla, retorna un
// Usuario vacío junto con un error descriptivo.
func NuevoUsuario(nombre, email, password string, planInicial suscripcion.Plan) (*Usuario, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if email == "" {
		return nil, errors.New("el email no puede estar vacío")
	}
	if len(password) < 4 {
		return nil, errors.New("la contraseña debe tener al menos 4 caracteres")
	}

	u := &Usuario{
		id:       contadorID,
		nombre:   nombre,
		email:    email,
		password: password,
		plan:     planInicial,
		activo:   true,
	}
	contadorID++
	return u, nil
}

// GetID retorna el identificador del usuario.
func (u *Usuario) GetID() int {
	return u.id
}

// GetNombre retorna el nombre del usuario.
func (u *Usuario) GetNombre() string {
	return u.nombre
}

// GetEmail retorna el correo del usuario.
func (u *Usuario) GetEmail() string {
	return u.email
}

// GetPlan retorna el plan de suscripción actual del usuario.
func (u *Usuario) GetPlan() suscripcion.Plan {
	return u.plan
}

// EstaActivo indica si la cuenta del usuario está activa.
func (u *Usuario) EstaActivo() bool {
	return u.activo
}

// Login valida las credenciales del usuario.
//
// MANEJO DE ERRORES: en lugar de retornar simplemente "true/false",
// retornamos un error descriptivo. Esto permite a quien llama saber
// EXACTAMENTE por qué falló el inicio de sesión.
func (u *Usuario) Login(emailIngresado, passwordIngresado string) error {
	if !u.activo {
		return errors.New("la cuenta de usuario está inactiva")
	}
	if u.email != emailIngresado {
		return errors.New("el correo ingresado no coincide")
	}
	if u.password != passwordIngresado {
		return errors.New("contraseña incorrecta")
	}
	return nil
}

// CambiarPassword permite actualizar la contraseña, pero exige
// conocer la contraseña actual.
//
// ENCAPSULACIÓN: esta es la ÚNICA forma de modificar u.password
// desde fuera del paquete.
func (u *Usuario) CambiarPassword(actual, nueva string) error {
	if u.password != actual {
		return errors.New("la contraseña actual no es correcta")
	}
	if len(nueva) < 4 {
		return errors.New("la nueva contraseña debe tener al menos 4 caracteres")
	}
	u.password = nueva
	return nil
}

// CambiarPlan asigna un nuevo plan de suscripción al usuario.
//
// Como "plan" es de tipo interfaz suscripcion.Plan, este método
// acepta CUALQUIER struct que implemente esa interfaz (PlanBasico,
// PlanEstandar, PlanPremium) — POLIMORFISMO.
func (u *Usuario) CambiarPlan(p suscripcion.Plan) {
	u.plan = p
}

// VerificarAcceso indica si el plan actual del usuario permite
// reproducir contenido con la calidad solicitada.
func (u *Usuario) VerificarAcceso(calidadSolicitada string) bool {
	return u.plan.CalidadMaxima() == calidadSolicitada || calidadSolicitada == "SD"
}

// Resumen retorna una descripción legible del usuario, útil para
// listados en consola.
func (u *Usuario) Resumen() string {
	estado := "activo"
	if !u.activo {
		estado = "inactivo"
	}
	return fmt.Sprintf("[%d] %s <%s> - Plan: %s (%s)",
		u.id, u.nombre, u.email, u.plan.Nombre(), estado)
}
