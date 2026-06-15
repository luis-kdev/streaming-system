// Package contenido gestiona el catálogo de contenido multimedia
// (películas y series) del sistema de streaming.
package contenido

import "fmt"

// Reproducible es la interfaz que deben cumplir todos los tipos de
// contenido que pueden reproducirse (Pelicula, Serie).
//
// Gracias a esta interfaz, el módulo "reproductor" puede trabajar
// con cualquier tipo de contenido sin conocer su tipo concreto.
// Esto es POLIMORFISMO: distintos structs, mismo "contrato".
type Reproducible interface {
	Reproducir() string
	Detalles() string
	IncrementarReproduccion()
	GetTitulo() string
	GetID() int
	GetReproducciones() int
}

// ContenidoBase contiene los atributos comunes a cualquier tipo de
// contenido multimedia.
//
// ENCAPSULACIÓN: los campos están en minúscula (privados), por lo que
// solo se pueden leer o modificar mediante los métodos públicos
// definidos abajo (GetTitulo, GetID, etc).
//
// HERENCIA: Pelicula y Serie van a EMBEBER este struct, heredando
// automáticamente estos campos y métodos.
type ContenidoBase struct {
	id             int
	titulo         string
	genero         string
	duracion       int // duración en minutos
	reproducciones int
}

// NuevoContenidoBase actúa como "constructor" de ContenidoBase.
func NuevoContenidoBase(id int, titulo, genero string, duracion int) ContenidoBase {
	return ContenidoBase{
		id:             id,
		titulo:         titulo,
		genero:         genero,
		duracion:       duracion,
		reproducciones: 0,
	}
}

// GetID retorna el identificador del contenido.
func (c ContenidoBase) GetID() int {
	return c.id
}

// GetTitulo retorna el título del contenido.
func (c ContenidoBase) GetTitulo() string {
	return c.titulo
}

// GetGenero retorna el género del contenido.
func (c ContenidoBase) GetGenero() string {
	return c.genero
}

// GetDuracion retorna la duración en minutos.
func (c ContenidoBase) GetDuracion() int {
	return c.duracion
}

// GetReproducciones retorna cuántas veces se ha reproducido el contenido.
func (c ContenidoBase) GetReproducciones() int {
	return c.reproducciones
}

// IncrementarReproduccion suma 1 al contador de reproducciones.
//
// Usamos un puntero (*ContenidoBase) como receptor porque queremos
// modificar el valor original almacenado en memoria, no una copia.
func (c *ContenidoBase) IncrementarReproduccion() {
	c.reproducciones++
}

// resumenBase construye un texto reutilizable con la info básica del
// contenido. Es una función privada (minúscula): solo Pelicula y Serie,
// al estar en el mismo paquete, pueden usarla.
func (c ContenidoBase) resumenBase() string {
	return fmt.Sprintf("[%d] %s (%s) - %d min - %d reproducciones",
		c.id, c.titulo, c.genero, c.duracion, c.reproducciones)
}
