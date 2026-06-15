// Package reportes genera estadísticas sobre el uso del sistema:
// contenido más visto y usuarios activos.
package reportes

import (
	"fmt"
	"sort"
	"strings"

	"streaming-system/contenido"
	"streaming-system/usuarios"
)

// Reportable es la interfaz que debe cumplir cualquier tipo de reporte.
//
// INTERFAZ: gracias a este "contrato", el sistema puede generar e
// imprimir cualquier reporte de la misma forma (MostrarReporte, abajo)
// sin importar su contenido interno. Si en el futuro se necesita un
// nuevo reporte (ej. ReporteIngresos), basta con crear un struct que
// implemente GenerarReporte() — PRINCIPIO ABIERTO/CERRADO.
type Reportable interface {
	GenerarReporte() string
}

// ReporteContenido genera un top del contenido más reproducido.
//
// ENCAPSULACIÓN: el slice "catalogo" es privado.
type ReporteContenido struct {
	catalogo []contenido.Reproducible
}

// NuevoReporteContenido crea un reporte a partir del catálogo actual.
func NuevoReporteContenido(catalogo []contenido.Reproducible) ReporteContenido {
	return ReporteContenido{catalogo: catalogo}
}

// GenerarReporte implementa la interfaz Reportable.
//
// POLIMORFISMO: "catalogo" mezcla *Pelicula y *Serie, pero ambos
// implementan GetTitulo() y GetReproducciones() (heredados de
// ContenidoBase), por lo que aquí se tratan de forma uniforme, sin
// necesidad de preguntar de qué tipo es cada uno.
func (r ReporteContenido) GenerarReporte() string {
	// Copiamos el slice para no alterar el orden original del catálogo.
	ordenado := make([]contenido.Reproducible, len(r.catalogo))
	copy(ordenado, r.catalogo)

	sort.Slice(ordenado, func(i, j int) bool {
		return ordenado[i].GetReproducciones() > ordenado[j].GetReproducciones()
	})

	var sb strings.Builder
	sb.WriteString("=== Top contenido más reproducido ===\n")

	limite := len(ordenado)
	if limite > 5 {
		limite = 5
	}
	for i := 0; i < limite; i++ {
		c := ordenado[i]
		sb.WriteString(fmt.Sprintf("%d. %s - %d reproducciones\n",
			i+1, c.GetTitulo(), c.GetReproducciones()))
	}
	return sb.String()
}

// ReporteUsuarios genera un listado de usuarios activos.
type ReporteUsuarios struct {
	lista []*usuarios.Usuario
}

// NuevoReporteUsuarios crea un reporte a partir de la lista de usuarios.
func NuevoReporteUsuarios(lista []*usuarios.Usuario) ReporteUsuarios {
	return ReporteUsuarios{lista: lista}
}

// GenerarReporte implementa la interfaz Reportable.
func (r ReporteUsuarios) GenerarReporte() string {
	var sb strings.Builder
	sb.WriteString("=== Usuarios activos ===\n")

	activos := 0
	for _, u := range r.lista {
		if u.EstaActivo() {
			activos++
			_, _ = sb.WriteString(u.Resumen() + "\n")
		}
	}
	sb.WriteString(fmt.Sprintf("Total activos: %d de %d\n", activos, len(r.lista)))
	return sb.String()
}

// MostrarReporte imprime en consola cualquier reporte que implemente
// la interfaz Reportable.
//
// POLIMORFISMO: esta función recibe tanto un ReporteContenido como un
// ReporteUsuarios (o cualquier futuro tipo de reporte) sin necesitar
// un if/switch para cada caso.
func MostrarReporte(r Reportable) {
	fmt.Println(r.GenerarReporte())
}
