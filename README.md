# Sistema de Gestión de Streaming

**Materia:** Programación Orientada a Objetos  
**Estudiante:** Luis Caizatoa | luis-kdev  
**Repositorio:** https://github.com/luis-kdev/streaming-system  
**Fecha:** 28 Junio 2026  

---

## Objetivo

Desarrollar un Sistema de Gestión de Streaming aplicando los principios
de Programación Orientada a Objetos en Go (Golang), evolucionando desde
una aplicación de consola hasta una API REST con servicios web.

---

## Tecnologías

- **Lenguaje:** Go 1.26
- **Paquetes:** net/http, encoding/json, sync, testing
- **Control de versiones:** Git + GitHub

---

## Estructura del proyecto

streaming-system/

├── api/

│   ├── server.go       ← 9 servicios web REST (Unidad 4)

│   └── server_test.go  ← 8 pruebas unitarias

├── contenido/

│   ├── base.go         ← interfaz Reproducible + ContenidoBase

│   ├── pelicula.go     ← struct Pelicula (herencia)

│   └── serie.go        ← struct Serie (herencia)

├── reportes/

│   └── reportes.go     ← interfaz Reportable

├── reproductor/

│   └── reproductor.go  ← struct Sesion

├── suscripcion/

│   └── suscripcion.go  ← interfaz Plan (3 implementaciones)

├── usuarios/

│   └── usuario.go      ← struct Usuario (encapsulación)

├── main.go             ← punto de entrada

└── go.mod

---

## Servicios Web (API REST)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/contenido | Listar catálogo completo |
| POST | /api/contenido/pelicula | Agregar película |
| POST | /api/contenido/serie | Agregar serie |
| GET | /api/usuarios | Listar usuarios |
| POST | /api/usuarios | Registrar usuario |
| POST | /api/login | Iniciar sesión |
| POST | /api/reproducir | Reproducir contenido |
| GET | /api/reportes | Ver reportes |
| GET | /api/planes | Ver planes de suscripción |

---

## Cómo ejecutar

```bash
git clone https://github.com/luis-kdev/streaming-system.git
cd streaming-system
go run .
```

El servidor inicia en: http://localhost:8080

## Cómo ejecutar los tests

```bash
go test ./api/... -v
```

---

## Conceptos POO aplicados

| Concepto | Dónde |
|----------|-------|
| Herencia (embedding) | Pelicula y Serie embeben ContenidoBase |
| Encapsulación | Campos privados + getters en todos los structs |
| Interfaces | Reproducible, Plan, Reportable |
| Polimorfismo | Reproducir(), MostrarReporte(), planes de suscripción |
| Concurrencia | sync.Mutex + goroutines automáticas en net/http |
| Manejo de errores | Patrón (resultado, error) en constructores |
| Serialización JSON | encoding/json en todos los servicios web |
