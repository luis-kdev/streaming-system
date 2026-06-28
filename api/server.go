// Package api implementa los servicios web REST del Sistema de Gestión
// de Streaming. Aplica los conceptos de la Unidad 4: concurrencia
// mediante goroutines, serialización JSON y servicios web RESTful.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"streaming-system/contenido"
	"streaming-system/reportes"
	"streaming-system/reproductor"
	"streaming-system/suscripcion"
	"streaming-system/usuarios"
)

// CONCURRENCIA: mu protege el acceso concurrente a los datos compartidos.
// Cuando dos peticiones llegan al mismo tiempo (dos goroutines), el Mutex
// asegura que solo una modifique los datos a la vez.
var (
	mu            sync.Mutex
	catalogo      []contenido.Reproducible
	listaUsuarios []*usuarios.Usuario
	sesiones      []*reproductor.Sesion
)

// RespuestaJSON es la estructura estándar de todas las respuestas.
// SERIALIZACIÓN JSON: las etiquetas `json:"..."` definen el nombre
// de cada campo en el JSON resultante.
type RespuestaJSON struct {
	OK      bool        `json:"ok"`
	Mensaje string      `json:"mensaje"`
	Datos   interface{} `json:"datos,omitempty"`
}

// ContenidoInfo estructura para serializar contenido a JSON.
type ContenidoInfo struct {
	ID             int    `json:"id"`
	Titulo         string `json:"titulo"`
	Detalles       string `json:"detalles"`
	Reproducciones int    `json:"reproducciones"`
}

// UsuarioInfo estructura para serializar usuarios a JSON.
type UsuarioInfo struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
	Plan   string `json:"plan"`
	Activo bool   `json:"activo"`
}

// PlanInfo estructura para serializar planes a JSON.
type PlanInfo struct {
	Nombre    string  `json:"nombre"`
	Precio    float64 `json:"precio"`
	Calidad   string  `json:"calidad_maxima"`
	Pantallas int     `json:"max_pantallas"`
}

// responder serializa la respuesta a JSON y la envía al cliente.
func responder(w http.ResponseWriter, codigo int, resp RespuestaJSON) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(resp)
}

// IniciarServidor configura las rutas y arranca el servidor HTTP.
// CONCURRENCIA: net/http lanza automáticamente una goroutine por cada
// petición entrante — concurrencia real sin código adicional.
func IniciarServidor(port string) {
	inicializarDatos()

	http.HandleFunc("/api/contenido", handlerListarContenido)
	http.HandleFunc("/api/contenido/pelicula", handlerAgregarPelicula)
	http.HandleFunc("/api/contenido/serie", handlerAgregarSerie)
	http.HandleFunc("/api/usuarios", handlerUsuarios)
	http.HandleFunc("/api/login", handlerLogin)
	http.HandleFunc("/api/reproducir", handlerReproducir)
	http.HandleFunc("/api/reportes", handlerReportes)
	http.HandleFunc("/api/planes", handlerPlanes)

	fmt.Printf("Servidor en http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}

// inicializarDatos carga datos de ejemplo al arrancar el servidor.
func inicializarDatos() {
	p1, _ := contenido.NuevaPelicula(1, "Interestelar", "Ciencia Ficción", "Christopher Nolan", 169)
	p2, _ := contenido.NuevaPelicula(2, "El Padrino", "Drama", "Francis Ford Coppola", 175)
	s1, _ := contenido.NuevaSerie(3, "Breaking Bad", "Drama", 47, 5, 62)
	s2, _ := contenido.NuevaSerie(4, "Stranger Things", "Misterio", 50, 4, 34)
	catalogo = append(catalogo, &p1, &p2, &s1, &s2)

	u1, _ := usuarios.NuevoUsuario("Luis Caizatoa", "luis@mail.com", "1234", suscripcion.PlanBasico{})
	u2, _ := usuarios.NuevoUsuario("Ana Torres", "ana@mail.com", "abcd", suscripcion.PlanPremium{})
	listaUsuarios = append(listaUsuarios, u1, u2)
}

// ── SERVICIO 1: GET /api/contenido ───────────────────────────────────────────
func handlerListarContenido(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	var items []ContenidoInfo
	for _, c := range catalogo {
		items = append(items, ContenidoInfo{
			ID:             c.GetID(),
			Titulo:         c.GetTitulo(),
			Detalles:       c.Detalles(),
			Reproducciones: c.GetReproducciones(),
		})
	}
	responder(w, 200, RespuestaJSON{OK: true, Mensaje: "Catálogo obtenido", Datos: items})
}

// ── SERVICIO 2: POST /api/contenido/pelicula ─────────────────────────────────
func handlerAgregarPelicula(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	var req struct {
		Titulo   string `json:"titulo"`
		Genero   string `json:"genero"`
		Director string `json:"director"`
		Duracion int    `json:"duracion"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: "JSON inválido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	id := len(catalogo) + 1
	p, err := contenido.NuevaPelicula(id, req.Titulo, req.Genero, req.Director, req.Duracion)
	if err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: err.Error()})
		return
	}
	catalogo = append(catalogo, &p)
	responder(w, 201, RespuestaJSON{OK: true, Mensaje: "Película agregada con ID " + strconv.Itoa(id)})
}

// ── SERVICIO 3: POST /api/contenido/serie ────────────────────────────────────
func handlerAgregarSerie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	var req struct {
		Titulo     string `json:"titulo"`
		Genero     string `json:"genero"`
		Duracion   int    `json:"duracion"`
		Temporadas int    `json:"temporadas"`
		Episodios  int    `json:"episodios"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: "JSON inválido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	id := len(catalogo) + 1
	s, err := contenido.NuevaSerie(id, req.Titulo, req.Genero, req.Duracion, req.Temporadas, req.Episodios)
	if err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: err.Error()})
		return
	}
	catalogo = append(catalogo, &s)
	responder(w, 201, RespuestaJSON{OK: true, Mensaje: "Serie agregada con ID " + strconv.Itoa(id)})
}

// ── SERVICIO 4 y 5: GET y POST /api/usuarios ─────────────────────────────────
func handlerUsuarios(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listarUsuarios(w, r)
	case http.MethodPost:
		registrarUsuario(w, r)
	default:
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
	}
}

func listarUsuarios(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var items []UsuarioInfo
	for _, u := range listaUsuarios {
		items = append(items, UsuarioInfo{
			ID:     u.GetID(),
			Nombre: u.GetNombre(),
			Email:  u.GetEmail(),
			Plan:   u.GetPlan().Nombre(),
			Activo: u.EstaActivo(),
		})
	}
	responder(w, 200, RespuestaJSON{OK: true, Mensaje: "Usuarios obtenidos", Datos: items})
}

func registrarUsuario(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nombre   string `json:"nombre"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Plan     string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: "JSON inválido"})
		return
	}
	var plan suscripcion.Plan
	switch req.Plan {
	case "estandar":
		plan = suscripcion.PlanEstandar{}
	case "premium":
		plan = suscripcion.PlanPremium{}
	default:
		plan = suscripcion.PlanBasico{}
	}
	u, err := usuarios.NuevoUsuario(req.Nombre, req.Email, req.Password, plan)
	if err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: err.Error()})
		return
	}
	mu.Lock()
	listaUsuarios = append(listaUsuarios, u)
	mu.Unlock()
	responder(w, 201, RespuestaJSON{OK: true, Mensaje: "Usuario registrado: " + req.Nombre})
}

// ── SERVICIO 6: POST /api/login ──────────────────────────────────────────────
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: "JSON inválido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for _, u := range listaUsuarios {
		if err := u.Login(req.Email, req.Password); err == nil {
			responder(w, 200, RespuestaJSON{
				OK:      true,
				Mensaje: "Login exitoso",
				Datos:   map[string]string{"nombre": u.GetNombre(), "plan": u.GetPlan().Nombre()},
			})
			return
		}
	}
	responder(w, 401, RespuestaJSON{OK: false, Mensaje: "Credenciales incorrectas"})
}

// ── SERVICIO 7: POST /api/reproducir ─────────────────────────────────────────
func handlerReproducir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	var req struct {
		UsuarioID   int `json:"usuario_id"`
		ContenidoID int `json:"contenido_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder(w, 400, RespuestaJSON{OK: false, Mensaje: "JSON inválido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	var c contenido.Reproducible
	for _, item := range catalogo {
		if item.GetID() == req.ContenidoID {
			c = item
			break
		}
	}
	if c == nil {
		responder(w, 404, RespuestaJSON{OK: false, Mensaje: "Contenido no encontrado"})
		return
	}
	sesion := reproductor.NuevaSesion(req.UsuarioID, c)
	sesion.Iniciar()
	sesiones = append(sesiones, sesion)
	responder(w, 200, RespuestaJSON{
		OK:      true,
		Mensaje: "Reproduciendo: " + c.GetTitulo(),
		Datos:   map[string]string{"estado": sesion.Resumen()},
	})
}

// ── SERVICIO 8: GET /api/reportes ────────────────────────────────────────────
func handlerReportes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	rContenido := reportes.NuevoReporteContenido(catalogo)
	rUsuarios := reportes.NuevoReporteUsuarios(listaUsuarios)
	responder(w, 200, RespuestaJSON{
		OK:      true,
		Mensaje: "Reportes generados",
		Datos: map[string]string{
			"contenido": rContenido.GenerarReporte(),
			"usuarios":  rUsuarios.GenerarReporte(),
		},
	})
}

// ── SERVICIO 9: GET /api/planes ──────────────────────────────────────────────
func handlerPlanes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responder(w, 405, RespuestaJSON{OK: false, Mensaje: "Método no permitido"})
		return
	}
	planes := []suscripcion.Plan{
		suscripcion.PlanBasico{},
		suscripcion.PlanEstandar{},
		suscripcion.PlanPremium{},
	}
	var items []PlanInfo
	for _, p := range planes {
		items = append(items, PlanInfo{
			Nombre:    p.Nombre(),
			Precio:    p.Precio(),
			Calidad:   p.CalidadMaxima(),
			Pantallas: p.MaxPantallas(),
		})
	}
	responder(w, 200, RespuestaJSON{OK: true, Mensaje: "Planes disponibles", Datos: items})
}
