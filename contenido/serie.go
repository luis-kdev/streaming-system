package contenido

import (
	"errors"
	"fmt"
)

// Serie representa una serie dentro del catálogo.
//
// HERENCIA: igual que Pelicula, Serie EMBEBE ContenidoBase y hereda
// sus campos (id, titulo, genero, duracion, reproducciones) y
// métodos (GetID, GetTitulo, etc).
type Serie struct {
	ContenidoBase
	temporadas int
	episodios  int
}

// NuevaSerie crea una nueva serie, validando los datos de entrada.
//
// MANEJO DE ERRORES: se valida que las temporadas y episodios sean
// valores positivos antes de crear el registro.
func NuevaSerie(id int, titulo, genero string, duracionPromedio, temporadas, episodios int) (Serie, error) {
	if titulo == "" {
		return Serie{}, errors.New("el título de la serie no puede estar vacío")
	}
	if temporadas <= 0 || episodios <= 0 {
		return Serie{}, errors.New("temporadas y episodios deben ser mayores a 0")
	}

	return Serie{
		ContenidoBase: NuevoContenidoBase(id, titulo, genero, duracionPromedio),
		temporadas:    temporadas,
		episodios:     episodios,
	}, nil
}

// GetTemporadas retorna el número de temporadas de la serie.
func (s Serie) GetTemporadas() int {
	return s.temporadas
}

// GetEpisodios retorna el número total de episodios de la serie.
func (s Serie) GetEpisodios() int {
	return s.episodios
}

// Reproducir marca la serie como reproducida una vez más.
//
// POLIMORFISMO: misma firma que Pelicula.Reproducir() (cumple la
// misma interfaz Reproducible), pero el mensaje es distinto porque
// una serie se reproduce "por episodios y temporadas".
func (s *Serie) Reproducir() string {
	s.IncrementarReproduccion()
	return fmt.Sprintf("Reproduciendo serie: %s (T%d, %d episodios)", s.GetTitulo(), s.temporadas, s.episodios)
}

// Detalles retorna información completa de la serie, combinando los
// datos heredados de ContenidoBase con los propios de Serie.
func (s Serie) Detalles() string {
	return fmt.Sprintf("%s | %d temporadas | %d episodios", s.resumenBase(), s.temporadas, s.episodios)
}
