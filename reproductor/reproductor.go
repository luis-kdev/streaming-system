// Package reproductor controla las sesiones de reproducción de
// contenido por parte de los usuarios.
package reproductor

import (
	"errors"
	"fmt"

	"streaming-system/contenido"
)

// Posibles estados de una sesión de reproducción.
const (
	EstadoReproduciendo = "reproduciendo"
	EstadoPausado       = "pausado"
	EstadoFinalizado    = "finalizado"
)

// Sesion representa una sesión de reproducción de un usuario sobre
// un contenido específico.
//
// ENCAPSULACIÓN: todos los campos son privados.
//
// RELACIÓN: "contenido" es de tipo interfaz contenido.Reproducible —
// puede ser *Pelicula o *Serie (POLIMORFISMO), sin que Sesion necesite
// saber cuál es.
type Sesion struct {
	usuarioID int
	contenido contenido.Reproducible
	progreso  int // segundos reproducidos
	estado    string
}

// NuevaSesion crea una sesión en estado "pausado", lista para iniciar.
func NuevaSesion(usuarioID int, c contenido.Reproducible) *Sesion {
	return &Sesion{
		usuarioID: usuarioID,
		contenido: c,
		progreso:  0,
		estado:    EstadoPausado,
	}
}

// Iniciar comienza (o reanuda) la reproducción.
//
// MANEJO DE ERRORES: no se puede iniciar una sesión ya finalizada.
func (s *Sesion) Iniciar() error {
	if s.estado == EstadoFinalizado {
		return errors.New("no se puede iniciar una sesión finalizada")
	}
	s.estado = EstadoReproduciendo
	// POLIMORFISMO: Reproducir() ejecuta la lógica correspondiente
	// según el tipo concreto (Pelicula o Serie).
	mensaje := s.contenido.Reproducir()
	fmt.Println(mensaje)
	return nil
}

// Pausar detiene temporalmente la reproducción.
func (s *Sesion) Pausar() error {
	if s.estado != EstadoReproduciendo {
		return errors.New("la sesión no se está reproduciendo")
	}
	s.estado = EstadoPausado
	return nil
}

// Reanudar continúa la reproducción desde donde se pausó.
func (s *Sesion) Reanudar() error {
	if s.estado != EstadoPausado {
		return errors.New("la sesión no está pausada")
	}
	return s.Iniciar()
}

// Finalizar marca la sesión como terminada.
func (s *Sesion) Finalizar() {
	s.estado = EstadoFinalizado
}

// AvanzarProgreso suma segundos al progreso de la sesión.
func (s *Sesion) AvanzarProgreso(segundos int) {
	s.progreso += segundos
}

// Resumen retorna una descripción legible de la sesión, incluyendo
// el título del contenido (vía Reproducible.GetTitulo()).
func (s *Sesion) Resumen() string {
	return fmt.Sprintf("Usuario #%d -> %s | Estado: %s | Progreso: %ds",
		s.usuarioID, s.contenido.GetTitulo(), s.estado, s.progreso)
}
