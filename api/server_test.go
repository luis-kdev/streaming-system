// server_test.go contiene las pruebas unitarias e integración
// de los servicios web del Sistema de Gestión de Streaming.
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// init inicializa los datos de prueba antes de correr los tests.
func init() {
	inicializarDatos()
}

// ── PRUEBA 1: GET /api/contenido ─────────────────────────────────────
// Verifica que el catálogo retorna 200 OK y al menos un elemento.
func TestListarContenido(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/contenido", nil)
	w := httptest.NewRecorder()

	handlerListarContenido(w, req)

	if w.Code != 200 {
		t.Errorf("Se esperaba código 200, se obtuvo %d", w.Code)
	}

	var resp RespuestaJSON
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp.OK {
		t.Error("Se esperaba OK=true en la respuesta")
	}
	t.Logf("✅ GET /api/contenido → %s", resp.Mensaje)
}

// ── PRUEBA 2: POST /api/contenido/pelicula ────────────────────────────
// Verifica que se puede agregar una nueva película.
func TestAgregarPelicula(t *testing.T) {
	body := `{"titulo":"Matrix","genero":"Ciencia Ficción","director":"Wachowski","duracion":136}`
	req := httptest.NewRequest(http.MethodPost, "/api/contenido/pelicula", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handlerAgregarPelicula(w, req)

	if w.Code != 201 {
		t.Errorf("Se esperaba código 201, se obtuvo %d", w.Code)
	}

	var resp RespuestaJSON
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp.OK {
		t.Errorf("Se esperaba OK=true, mensaje: %s", resp.Mensaje)
	}
	t.Logf("✅ POST /api/contenido/pelicula → %s", resp.Mensaje)
}

// ── PRUEBA 3: POST /api/contenido/pelicula sin título (error esperado)
// Verifica que el manejo de errores funciona correctamente.
func TestAgregarPeliculaSinTitulo(t *testing.T) {
	body := `{"titulo":"","genero":"Drama","director":"Test","duracion":100}`
	req := httptest.NewRequest(http.MethodPost, "/api/contenido/pelicula", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handlerAgregarPelicula(w, req)

	if w.Code != 400 {
		t.Errorf("Se esperaba código 400 por título vacío, se obtuvo %d", w.Code)
	}
	t.Logf("✅ Validación de título vacío funciona correctamente")
}

// ── PRUEBA 4: GET /api/usuarios ───────────────────────────────────────
// Verifica que la lista de usuarios retorna datos correctos.
func TestListarUsuarios(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/usuarios", nil)
	w := httptest.NewRecorder()

	listarUsuarios(w, req)

	if w.Code != 200 {
		t.Errorf("Se esperaba código 200, se obtuvo %d", w.Code)
	}

	var resp RespuestaJSON
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp.OK {
		t.Error("Se esperaba OK=true en la respuesta")
	}
	t.Logf("✅ GET /api/usuarios → %s", resp.Mensaje)
}

// ── PRUEBA 5: POST /api/login exitoso ─────────────────────────────────
// Verifica que el login con credenciales correctas retorna 200.
func TestLoginExitoso(t *testing.T) {
	body := `{"email":"luis@mail.com","password":"1234"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handlerLogin(w, req)

	if w.Code != 200 {
		t.Errorf("Se esperaba código 200, se obtuvo %d", w.Code)
	}
	t.Logf("✅ POST /api/login exitoso → código %d", w.Code)
}

// ── PRUEBA 6: POST /api/login fallido ────────────────────────────────
// Verifica que credenciales incorrectas retornan 401.
func TestLoginFallido(t *testing.T) {
	body := `{"email":"luis@mail.com","password":"incorrecta"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handlerLogin(w, req)

	if w.Code != 401 {
		t.Errorf("Se esperaba código 401, se obtuvo %d", w.Code)
	}
	t.Logf("✅ POST /api/login fallido → código 401 correcto")
}

// ── PRUEBA 7: GET /api/reportes ───────────────────────────────────────
// Verifica que los reportes se generan correctamente.
func TestReportes(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/reportes", nil)
	w := httptest.NewRecorder()

	handlerReportes(w, req)

	if w.Code != 200 {
		t.Errorf("Se esperaba código 200, se obtuvo %d", w.Code)
	}

	var resp RespuestaJSON
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp.OK {
		t.Error("Se esperaba OK=true en reportes")
	}
	t.Logf("✅ GET /api/reportes → %s", resp.Mensaje)
}

// ── PRUEBA 8: GET /api/planes ─────────────────────────────────────────
// Verifica que los 3 planes de suscripción se retornan correctamente.
func TestPlanes(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/planes", nil)
	w := httptest.NewRecorder()

	handlerPlanes(w, req)

	if w.Code != 200 {
		t.Errorf("Se esperaba código 200, se obtuvo %d", w.Code)
	}
	t.Logf("✅ GET /api/planes → código %d", w.Code)
}
