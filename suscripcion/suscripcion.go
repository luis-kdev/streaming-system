// Package suscripcion define los distintos planes de suscripción
// disponibles en el sistema de streaming.
package suscripcion

import "fmt"

// Plan es la interfaz que debe cumplir cualquier plan de suscripción.
//
// INTERFAZ + POLIMORFISMO: PlanBasico, PlanEstandar y PlanPremium son
// tipos completamente independientes entre sí (NO comparten un struct
// base como en el módulo "contenido"), pero los tres implementan este
// mismo "contrato". Así, el módulo usuarios puede manejar cualquier
// Plan sin importar cuál sea su tipo concreto.
type Plan interface {
	Nombre() string
	Precio() float64
	CalidadMaxima() string
	MaxPantallas() int
	Resumen() string
}

// PlanBasico representa el plan más económico del sistema.
type PlanBasico struct{}

func (p PlanBasico) Nombre() string        { return "Básico" }
func (p PlanBasico) Precio() float64       { return 4.99 }
func (p PlanBasico) CalidadMaxima() string { return "SD" }
func (p PlanBasico) MaxPantallas() int     { return 1 }

// Resumen retorna una descripción legible del plan. Cada plan
// implementa esta misma firma, pero con sus propios valores.
func (p PlanBasico) Resumen() string {
	return fmt.Sprintf("%s - $%.2f/mes - Calidad %s - %d pantalla(s)",
		p.Nombre(), p.Precio(), p.CalidadMaxima(), p.MaxPantallas())
}

// PlanEstandar representa el plan intermedio.
type PlanEstandar struct{}

func (p PlanEstandar) Nombre() string        { return "Estándar" }
func (p PlanEstandar) Precio() float64       { return 8.99 }
func (p PlanEstandar) CalidadMaxima() string { return "HD" }
func (p PlanEstandar) MaxPantallas() int     { return 2 }

func (p PlanEstandar) Resumen() string {
	return fmt.Sprintf("%s - $%.2f/mes - Calidad %s - %d pantalla(s)",
		p.Nombre(), p.Precio(), p.CalidadMaxima(), p.MaxPantallas())
}

// PlanPremium representa el plan más completo del sistema.
type PlanPremium struct{}

func (p PlanPremium) Nombre() string        { return "Premium" }
func (p PlanPremium) Precio() float64       { return 12.99 }
func (p PlanPremium) CalidadMaxima() string { return "4K" }
func (p PlanPremium) MaxPantallas() int     { return 4 }

func (p PlanPremium) Resumen() string {
	return fmt.Sprintf("%s - $%.2f/mes - Calidad %s - %d pantalla(s)",
		p.Nombre(), p.Precio(), p.CalidadMaxima(), p.MaxPantallas())
}
