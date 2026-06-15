// main.go conecta los 5 módulos del sistema (contenido, suscripcion,
// usuarios, reproductor, reportes) en un programa ejecutable con menú
// interactivo en consola.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"streaming-system/contenido"
	"streaming-system/reportes"
	"streaming-system/reproductor"
	"streaming-system/suscripcion"
	"streaming-system/usuarios"
)

// catalogo almacena todo el contenido disponible.
//
// POLIMORFISMO: es un slice de la interfaz contenido.Reproducible,
// por lo que puede contener *Pelicula y *Serie mezclados.
var catalogo []contenido.Reproducible

// listaUsuarios almacena todos los usuarios registrados.
var listaUsuarios []*usuarios.Usuario

// sesionActual es la sesión de reproducción activa (si existe).
var sesionActual *reproductor.Sesion

// usuarioActual es el usuario que inició sesión (si existe).
var usuarioActual *usuarios.Usuario

// lector permite leer la entrada del usuario desde la consola.
var lector = bufio.NewReader(os.Stdin)

func main() {
	inicializarDatos()

	for {
		mostrarMenu()
		switch leerLinea() {
		case "1":
			listarCatalogo()
		case "2":
			iniciarSesion()
		case "3":
			reproducirContenido()
		case "4":
			pausarOReanudar()
		case "5":
			verReportes()
		case "0":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

// inicializarDatos carga el catálogo y los usuarios de ejemplo.
//
// MANEJO DE ERRORES: si algún constructor (NuevaPelicula, NuevaSerie,
// NuevoUsuario) retorna error por datos inválidos, se imprime el
// error y ese elemento simplemente no se agrega.
func inicializarDatos() {
	p1, err := contenido.NuevaPelicula(1, "Interestelar", "Ciencia Ficción", "Christopher Nolan", 169)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		catalogo = append(catalogo, &p1)
	}

	p2, err := contenido.NuevaPelicula(2, "El Padrino", "Drama", "Francis Ford Coppola", 175)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		catalogo = append(catalogo, &p2)
	}

	s1, err := contenido.NuevaSerie(3, "Breaking Bad", "Drama", 47, 5, 62)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		catalogo = append(catalogo, &s1)
	}

	s2, err := contenido.NuevaSerie(4, "Stranger Things", "Misterio", 50, 4, 34)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		catalogo = append(catalogo, &s2)
	}

	// POLIMORFISMO: cada usuario recibe un tipo distinto de Plan
	// (PlanBasico, PlanPremium), pero el campo "plan" en Usuario es
	// del tipo interfaz suscripcion.Plan para ambos casos.
	u1, err := usuarios.NuevoUsuario("Luis Caizatoa", "luis@mail.com", "1234", suscripcion.PlanBasico{})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		listaUsuarios = append(listaUsuarios, u1)
	}

	u2, err := usuarios.NuevoUsuario("Ana Torres", "ana@mail.com", "abcd", suscripcion.PlanPremium{})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		listaUsuarios = append(listaUsuarios, u2)
	}
}

// mostrarMenu imprime las opciones disponibles en consola.
func mostrarMenu() {
	fmt.Println()
	fmt.Println("========== SISTEMA DE GESTIÓN DE STREAMING ==========")
	if usuarioActual != nil {
		fmt.Println("Usuario:", usuarioActual.Resumen())
	} else {
		fmt.Println("Usuario: (no ha iniciado sesión)")
	}
	fmt.Println("1. Ver catálogo")
	fmt.Println("2. Iniciar sesión")
	fmt.Println("3. Reproducir contenido")
	fmt.Println("4. Pausar / Reanudar sesión")
	fmt.Println("5. Ver reportes")
	fmt.Println("0. Salir")
	fmt.Print("Elige una opción: ")
}

// leerLinea lee una línea de texto desde la consola, sin el salto de línea.
func leerLinea() string {
	texto, _ := lector.ReadString('\n')
	return strings.TrimSpace(texto)
}

// listarCatalogo imprime el catálogo completo.
//
// POLIMORFISMO: Detalles() se llama sobre cada elemento sin saber si
// es *Pelicula o *Serie.
func listarCatalogo() {
	fmt.Println("\n--- Catálogo disponible ---")
	for _, c := range catalogo {
		fmt.Println(c.Detalles())
	}
}

// iniciarSesion solicita email y contraseña, y valida con Login().
//
// MANEJO DE ERRORES: Login() retorna un error específico; aquí solo
// mostramos un mensaje genérico para no revelar cuál dato fue
// incorrecto (buena práctica de seguridad).
func iniciarSesion() {
	fmt.Print("Email: ")
	email := leerLinea()
	fmt.Print("Contraseña: ")
	password := leerLinea()

	for _, u := range listaUsuarios {
		if err := u.Login(email, password); err == nil {
			usuarioActual = u
			fmt.Println("Sesión iniciada como", u.GetNombre())
			return
		}
	}
	fmt.Println("Error: credenciales incorrectas o usuario no encontrado")
}

// reproducirContenido crea una Sesion para el usuario actual sobre el
// contenido elegido.
func reproducirContenido() {
	if usuarioActual == nil {
		fmt.Println("Error: primero debes iniciar sesión (opción 2)")
		return
	}

	listarCatalogo()
	fmt.Print("ID del contenido a reproducir: ")
	id, err := strconv.Atoi(leerLinea())
	if err != nil {
		fmt.Println("Error: ID inválido")
		return
	}

	for _, c := range catalogo {
		if c.GetID() == id {
			sesionActual = reproductor.NuevaSesion(usuarioActual.GetID(), c)
			if err := sesionActual.Iniciar(); err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("Error: no se encontró contenido con ese ID")
}

// pausarOReanudar controla el estado de la sesión activa.
func pausarOReanudar() {
	if sesionActual == nil {
		fmt.Println("Error: no hay ninguna sesión activa")
		return
	}

	fmt.Println(sesionActual.Resumen())
	fmt.Println("1. Pausar  2. Reanudar")
	fmt.Print("Elige: ")

	var err error
	switch leerLinea() {
	case "1":
		err = sesionActual.Pausar()
	case "2":
		err = sesionActual.Reanudar()
	default:
		fmt.Println("Opción no válida")
		return
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}

// verReportes muestra el top de contenido y los usuarios activos.
//
// POLIMORFISMO: MostrarReporte() recibe primero un ReporteContenido y
// luego un ReporteUsuarios — mismo método, distinto tipo concreto.
func verReportes() {
	fmt.Println()
	reportes.MostrarReporte(reportes.NuevoReporteContenido(catalogo))
	fmt.Println()
	reportes.MostrarReporte(reportes.NuevoReporteUsuarios(listaUsuarios))
}
