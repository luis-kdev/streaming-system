package contenido

import (
	"errors"
	"fmt"
)

// Pelicula representa una película dentro del catálogo.
//
// HERENCIA: al EMBEBER ContenidoBase (sin nombre de campo), Pelicula
// "hereda" automáticamente los campos id, titulo, genero, duracion,
// reproducciones, y los métodos GetID(), GetTitulo(), etc.
type Pelicula struct {
	ContenidoBase
	director string
}

// NuevaPelicula crea una nueva película, validando los datos de entrada.
//
// MANEJO DE ERRORES: si algún dato no es válido, se retorna una
// Pelicula vacía junto con un error explicativo, en vez de crear
// un registro inválido en el sistema.
func NuevaPelicula(id int, titulo, genero, director string, duracion int) (Pelicula, error) {
	if titulo == "" {
		return Pelicula{}, errors.New("el título de la película no puede estar vacío")
	}
	if duracion <= 0 {
		return Pelicula{}, errors.New("la duración de la película debe ser mayor a 0")
	}

	return Pelicula{
		ContenidoBase: NuevoContenidoBase(id, titulo, genero, duracion),
		director:      director,
	}, nil
}

// GetDirector retorna el director de la película.
func (p Pelicula) GetDirector() string {
	return p.director
}

// Reproducir marca la película como reproducida una vez más y
// retorna un mensaje descriptivo.
//
// POLIMORFISMO: este método cumple el contrato de la interfaz
// Reproducible, pero su implementación es distinta a la de Serie
// (ver serie.go).
func (p *Pelicula) Reproducir() string {
	p.IncrementarReproduccion()
	return fmt.Sprintf("Reproduciendo película: %s", p.GetTitulo())
}

// Detalles retorna información completa de la película, combinando
// los datos heredados de ContenidoBase (resumenBase) con el dato
// propio de Pelicula (director).
func (p Pelicula) Detalles() string {
	return fmt.Sprintf("%s | Director: %s", p.resumenBase(), p.director)
}
