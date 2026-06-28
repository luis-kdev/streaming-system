// main.go — punto de entrada del Sistema de Gestión de Streaming.
// En esta versión final integra las 4 unidades:
// - Unidad 1: paquetes y funciones
// - Unidad 2: structs y arrays/slices
// - Unidad 3: encapsulación e interfaces
// - Unidad 4: servicios web con concurrencia y JSON
package main

import (
	"fmt"
	"streaming-system/api"
)

func main() {
	fmt.Println("=== SISTEMA DE GESTIÓN DE STREAMING ===")
	fmt.Println("Iniciando servidor de servicios web...")
	fmt.Println()

	// Inicia el servidor HTTP en el puerto 8080.
	// CONCURRENCIA: cada petición web se maneja en una goroutine
	// separada automáticamente por el paquete net/http de Go.
	api.IniciarServidor(":8080")
}
