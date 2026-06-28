// Package api - frontend.go sirve la interfaz web del sistema de streaming.
// El servidor Go entrega las páginas HTML directamente al navegador,
// convirtiendo el sistema en una aplicación web completa.
package api

import (
	"html/template"
	"net/http"
)

// paginaHTML es la interfaz web completa del sistema de streaming.
// Usa Bootstrap para el diseño visual y JavaScript (fetch API) para
// comunicarse con los servicios web REST que ya implementamos.
const paginaHTML = `
<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>StreamGo — Sistema de Gestión de Streaming</title>
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
<link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css" rel="stylesheet">
<style>
  body { background: #0f0f1a; color: #e0e0e0; font-family: 'Segoe UI', sans-serif; }
  .navbar { background: linear-gradient(135deg, #1e2761, #5340bb); }
  .card { background: #1a1a2e; border: 1px solid #2a2a4a; border-radius: 12px; }
  .card-header { background: linear-gradient(135deg, #5340bb, #0f6e56); border-radius: 12px 12px 0 0 !important; }
  .btn-primary { background: #5340bb; border-color: #5340bb; }
  .btn-primary:hover { background: #3d2d99; }
  .btn-success { background: #0f6e56; border-color: #0f6e56; }
  .table { color: #e0e0e0; }
  .table-dark { background: #12122a; }
  .badge-basico   { background: #0f6e56; }
  .badge-estandar { background: #ba7517; }
  .badge-premium  { background: #5340bb; }
  .hero { background: linear-gradient(135deg, #1e2761 0%, #0f0f1a 100%); 
          padding: 60px 0; text-align: center; }
  .hero h1 { font-size: 3rem; font-weight: bold; color: #cadcfc; }
  .hero p  { font-size: 1.2rem; color: #8a9ab5; }
  .section-title { color: #cadcfc; border-bottom: 2px solid #5340bb; padding-bottom: 8px; }
  #alertBox { display: none; }
  .content-card { transition: transform 0.2s; cursor: pointer; }
  .content-card:hover { transform: translateY(-4px); }
  .plan-card { border: 2px solid transparent; transition: all 0.3s; }
  .plan-card:hover { border-color: #5340bb; transform: translateY(-4px); }
  .nav-tabs .nav-link { color: #8a9ab5; }
  .nav-tabs .nav-link.active { color: #cadcfc; background: #1a1a2e; border-color: #5340bb; }
</style>
</head>
<body>

<!-- NAVBAR -->
<nav class="navbar navbar-dark px-4 py-3">
  <span class="navbar-brand fw-bold fs-4">
    <i class="bi bi-play-circle-fill me-2" style="color:#cadcfc"></i>StreamGo
  </span>
  <div id="userBadge" class="d-none">
    <span class="badge bg-light text-dark me-2" id="userNameBadge"></span>
    <span class="badge" id="userPlanBadge" style="background:#5340bb"></span>
    <button class="btn btn-sm btn-outline-light ms-2" onclick="cerrarSesion()">
      <i class="bi bi-box-arrow-right"></i> Salir
    </button>
  </div>
</nav>

<!-- HERO -->
<div class="hero">
  <h1><i class="bi bi-collection-play-fill"></i> Sistema de Gestión de Streaming</h1>
  <p>Aplicativo web desarrollado con Go (Golang) · POO · API REST · Concurrencia</p>
  <div class="mt-3">
    <span class="badge me-2" style="background:#5340bb;padding:8px 16px">Go 1.26</span>
    <span class="badge me-2" style="background:#0f6e56;padding:8px 16px">REST API</span>
    <span class="badge me-2" style="background:#ba7517;padding:8px 16px">JSON</span>
    <span class="badge" style="background:#993c1d;padding:8px 16px">Concurrencia</span>
  </div>
</div>

<!-- ALERTA -->
<div class="container mt-3">
  <div id="alertBox" class="alert alert-dismissible" role="alert">
    <span id="alertMsg"></span>
    <button type="button" class="btn-close" onclick="ocultarAlerta()"></button>
  </div>
</div>

<!-- TABS PRINCIPALES -->
<div class="container mt-4">
  <ul class="nav nav-tabs mb-4" id="mainTabs">
    <li class="nav-item">
      <a class="nav-link active" onclick="mostrarTab('catalogo')" href="#">
        <i class="bi bi-film"></i> Catálogo
      </a>
    </li>
    <li class="nav-item">
      <a class="nav-link" onclick="mostrarTab('planes')" href="#">
        <i class="bi bi-credit-card"></i> Planes
      </a>
    </li>
    <li class="nav-item">
      <a class="nav-link" onclick="mostrarTab('usuarios')" href="#">
        <i class="bi bi-people"></i> Usuarios
      </a>
    </li>
    <li class="nav-item">
      <a class="nav-link" onclick="mostrarTab('reportes')" href="#">
        <i class="bi bi-bar-chart"></i> Reportes
      </a>
    </li>
    <li class="nav-item">
      <a class="nav-link" onclick="mostrarTab('login')" href="#">
        <i class="bi bi-person-circle"></i> Login
      </a>
    </li>
  </ul>

  <!-- TAB: CATÁLOGO -->
  <div id="tab-catalogo">
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="section-title"><i class="bi bi-collection-play"></i> Catálogo de Contenido</h4>
      <button class="btn btn-success btn-sm" onclick="mostrarFormContenido()">
        <i class="bi bi-plus-circle"></i> Agregar Contenido
      </button>
    </div>

    <!-- Formulario agregar contenido -->
    <div id="formContenido" class="card p-3 mb-4" style="display:none">
      <h6 class="text-white mb-3"><i class="bi bi-plus"></i> Nuevo Contenido</h6>
      <div class="row g-2 mb-2">
        <div class="col-md-3">
          <select id="tipoContenido" class="form-select bg-dark text-white border-secondary" onchange="toggleFormContenido()">
            <option value="pelicula">🎬 Película</option>
            <option value="serie">📺 Serie</option>
          </select>
        </div>
        <div class="col-md-3"><input id="cTitulo"    class="form-control bg-dark text-white border-secondary" placeholder="Título *"></div>
        <div class="col-md-3"><input id="cGenero"    class="form-control bg-dark text-white border-secondary" placeholder="Género"></div>
        <div class="col-md-3"><input id="cDuracion"  class="form-control bg-dark text-white border-secondary" placeholder="Duración (min)" type="number"></div>
      </div>
      <div id="camposPelicula" class="row g-2 mb-2">
        <div class="col-md-4"><input id="cDirector" class="form-control bg-dark text-white border-secondary" placeholder="Director"></div>
      </div>
      <div id="camposSerie" class="row g-2 mb-2" style="display:none">
        <div class="col-md-3"><input id="cTemporadas" class="form-control bg-dark text-white border-secondary" placeholder="Temporadas" type="number"></div>
        <div class="col-md-3"><input id="cEpisodios"  class="form-control bg-dark text-white border-secondary" placeholder="Episodios"  type="number"></div>
      </div>
      <div class="d-flex gap-2">
        <button class="btn btn-primary btn-sm" onclick="agregarContenido()"><i class="bi bi-save"></i> Guardar</button>
        <button class="btn btn-secondary btn-sm" onclick="document.getElementById('formContenido').style.display='none'">Cancelar</button>
      </div>
    </div>

    <div id="gridCatalogo" class="row g-3"></div>
  </div>

  <!-- TAB: PLANES -->
  <div id="tab-planes" style="display:none">
    <h4 class="section-title mb-4"><i class="bi bi-credit-card"></i> Planes de Suscripción</h4>
    <div id="gridPlanes" class="row g-4 justify-content-center"></div>
  </div>

  <!-- TAB: USUARIOS -->
  <div id="tab-usuarios" style="display:none">
    <div class="d-flex justify-content-between align-items-center mb-3">
      <h4 class="section-title"><i class="bi bi-people"></i> Usuarios del Sistema</h4>
      <button class="btn btn-success btn-sm" onclick="mostrarFormUsuario()">
        <i class="bi bi-person-plus"></i> Nuevo Usuario
      </button>
    </div>

    <div id="formUsuario" class="card p-3 mb-4" style="display:none">
      <h6 class="text-white mb-3"><i class="bi bi-person-plus"></i> Registrar Usuario</h6>
      <div class="row g-2">
        <div class="col-md-3"><input id="uNombre"   class="form-control bg-dark text-white border-secondary" placeholder="Nombre *"></div>
        <div class="col-md-3"><input id="uEmail"    class="form-control bg-dark text-white border-secondary" placeholder="Email *" type="email"></div>
        <div class="col-md-2"><input id="uPassword" class="form-control bg-dark text-white border-secondary" placeholder="Contraseña *" type="password"></div>
        <div class="col-md-2">
          <select id="uPlan" class="form-select bg-dark text-white border-secondary">
            <option value="basico">Básico</option>
            <option value="estandar">Estándar</option>
            <option value="premium">Premium</option>
          </select>
        </div>
        <div class="col-md-2 d-flex gap-2">
          <button class="btn btn-primary btn-sm" onclick="registrarUsuario()"><i class="bi bi-save"></i> Guardar</button>
          <button class="btn btn-secondary btn-sm" onclick="document.getElementById('formUsuario').style.display='none'">X</button>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-body p-0">
        <table class="table table-dark table-hover mb-0">
          <thead><tr>
            <th>#</th><th>Nombre</th><th>Email</th><th>Plan</th><th>Estado</th>
          </tr></thead>
          <tbody id="tablaUsuarios"></tbody>
        </table>
      </div>
    </div>
  </div>

  <!-- TAB: REPORTES -->
  <div id="tab-reportes" style="display:none">
    <h4 class="section-title mb-4"><i class="bi bi-bar-chart"></i> Reportes del Sistema</h4>
    <div class="row g-4">
      <div class="col-md-6">
        <div class="card h-100">
          <div class="card-header"><i class="bi bi-trophy"></i> Top Contenido más Reproducido</div>
          <div class="card-body"><pre id="reporteContenido" style="color:#cadcfc;white-space:pre-wrap;font-size:0.9rem"></pre></div>
        </div>
      </div>
      <div class="col-md-6">
        <div class="card h-100">
          <div class="card-header"><i class="bi bi-people"></i> Usuarios Activos</div>
          <div class="card-body"><pre id="reporteUsuarios" style="color:#cadcfc;white-space:pre-wrap;font-size:0.9rem"></pre></div>
        </div>
      </div>
    </div>
  </div>

  <!-- TAB: LOGIN -->
  <div id="tab-login" style="display:none">
    <div class="row justify-content-center">
      <div class="col-md-5">
        <div class="card">
          <div class="card-header text-center py-3">
            <h5 class="mb-0"><i class="bi bi-person-circle"></i> Iniciar Sesión</h5>
          </div>
          <div class="card-body p-4">
            <div class="mb-3">
              <label class="form-label text-secondary">Email</label>
              <input id="loginEmail" class="form-control bg-dark text-white border-secondary" 
                     type="email" placeholder="luis@mail.com" value="luis@mail.com">
            </div>
            <div class="mb-4">
              <label class="form-label text-secondary">Contraseña</label>
              <input id="loginPassword" class="form-control bg-dark text-white border-secondary" 
                     type="password" placeholder="••••" value="1234">
            </div>
            <button class="btn btn-primary w-100" onclick="iniciarSesion()">
              <i class="bi bi-box-arrow-in-right"></i> Ingresar
            </button>
            <div class="mt-3 text-center text-secondary small">
              Usuario de prueba: luis@mail.com / 1234
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="mt-5 pb-4 text-center text-secondary small">
    StreamGo · Sistema de Gestión de Streaming · Go (Golang) · luis-kdev
    · <a href="https://github.com/luis-kdev/streaming-system" class="text-secondary">GitHub</a>
  </div>
</div>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
<script>
const API = '';
let usuarioActual = null;

// ── TABS ──────────────────────────────────────────────────────────────────────
function mostrarTab(tab) {
  ['catalogo','planes','usuarios','reportes','login'].forEach(t => {
    document.getElementById('tab-'+t).style.display = t===tab ? 'block' : 'none';
  });
  document.querySelectorAll('.nav-link').forEach(l => l.classList.remove('active'));
  event.target.classList.add('active');
  if (tab==='catalogo')  cargarCatalogo();
  if (tab==='planes')    cargarPlanes();
  if (tab==='usuarios')  cargarUsuarios();
  if (tab==='reportes')  cargarReportes();
}

// ── ALERTAS ──────────────────────────────────────────────────────────────────
function mostrarAlerta(msg, tipo='success') {
  const box = document.getElementById('alertBox');
  box.className = 'alert alert-'+tipo+' alert-dismissible';
  document.getElementById('alertMsg').textContent = msg;
  box.style.display = 'block';
  setTimeout(ocultarAlerta, 4000);
}
function ocultarAlerta() { document.getElementById('alertBox').style.display='none'; }

// ── CATÁLOGO ──────────────────────────────────────────────────────────────────
async function cargarCatalogo() {
  const r = await fetch(API+'/api/contenido');
  const data = await r.json();
  const grid = document.getElementById('gridCatalogo');
  if (!data.datos || data.datos.length===0) {
    grid.innerHTML = '<div class="col text-center text-secondary py-5"><i class="bi bi-film fs-1"></i><p class="mt-2">No hay contenido</p></div>';
    return;
  }
  grid.innerHTML = data.datos.map(c => {
    const esFilm = c.detalles.includes('Director');
    const icono  = esFilm ? 'bi-camera-reels' : 'bi-tv';
    const tipo   = esFilm ? '🎬 Película' : '📺 Serie';
    const color  = esFilm ? '#5340bb' : '#0f6e56';
    return ` + "`" + `
      <div class="col-md-3 col-sm-6">
        <div class="card content-card h-100" onclick="reproducir(${c.id}, '${c.titulo}')">
          <div class="card-body text-center py-4">
            <div class="mb-3" style="font-size:3rem"><i class="bi ${icono}" style="color:${color}"></i></div>
            <h6 class="fw-bold text-white">${c.titulo}</h6>
            <p class="text-secondary small mb-2">${c.detalles.split('|')[0].trim()}</p>
            <span class="badge" style="background:${color}">${tipo}</span>
            <div class="mt-2 text-secondary small">
              <i class="bi bi-play-circle"></i> ${c.reproducciones} reproducciones
            </div>
          </div>
          <div class="card-footer text-center py-2 border-0" style="background:rgba(255,255,255,0.05)">
            <small class="text-secondary"><i class="bi bi-hand-index"></i> Clic para reproducir</small>
          </div>
        </div>
      </div>` + "`" + `;
  }).join('');
}

async function reproducir(id, titulo) {
  if (!usuarioActual) {
    mostrarAlerta('Debes iniciar sesión primero (pestaña Login)', 'warning');
    return;
  }
  const r = await fetch(API+'/api/reproducir', {
    method: 'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify({ usuario_id: usuarioActual.id || 1, contenido_id: id })
  });
  const data = await r.json();
  if (data.ok) {
    mostrarAlerta('▶ Reproduciendo: ' + titulo, 'success');
    setTimeout(cargarCatalogo, 500);
  } else {
    mostrarAlerta(data.mensaje, 'danger');
  }
}

function mostrarFormContenido() {
  const f = document.getElementById('formContenido');
  f.style.display = f.style.display==='none' ? 'block' : 'none';
}

function toggleFormContenido() {
  const tipo = document.getElementById('tipoContenido').value;
  document.getElementById('camposPelicula').style.display = tipo==='pelicula' ? 'flex' : 'none';
  document.getElementById('camposSerie').style.display    = tipo==='serie'    ? 'flex' : 'none';
}

async function agregarContenido() {
  const tipo = document.getElementById('tipoContenido').value;
  let body, url;
  if (tipo==='pelicula') {
    url = '/api/contenido/pelicula';
    body = { titulo: document.getElementById('cTitulo').value,
             genero: document.getElementById('cGenero').value,
             director: document.getElementById('cDirector').value,
             duracion: parseInt(document.getElementById('cDuracion').value)||0 };
  } else {
    url = '/api/contenido/serie';
    body = { titulo: document.getElementById('cTitulo').value,
             genero: document.getElementById('cGenero').value,
             duracion: parseInt(document.getElementById('cDuracion').value)||0,
             temporadas: parseInt(document.getElementById('cTemporadas').value)||0,
             episodios:  parseInt(document.getElementById('cEpisodios').value)||0 };
  }
  const r = await fetch(API+url, { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(body) });
  const data = await r.json();
  mostrarAlerta(data.mensaje, data.ok ? 'success' : 'danger');
  if (data.ok) { document.getElementById('formContenido').style.display='none'; cargarCatalogo(); }
}

// ── PLANES ────────────────────────────────────────────────────────────────────
async function cargarPlanes() {
  const r = await fetch(API+'/api/planes');
  const data = await r.json();
  const colores = ['#0f6e56','#5340bb','#ba7517'];
  const iconos  = ['bi-wifi','bi-wifi-2','bi-stars'];
  document.getElementById('gridPlanes').innerHTML = (data.datos||[]).map((p,i) => ` + "`" + `
    <div class="col-md-3">
      <div class="card plan-card text-center p-4 h-100">
        <div class="mb-3" style="color:${colores[i]};font-size:2.5rem">
          <i class="bi ${iconos[i]}"></i>
        </div>
        <h5 class="fw-bold text-white">${p.nombre}</h5>
        <div class="my-3">
          <span style="font-size:2rem;color:${colores[i]};font-weight:bold">$${p.precio}</span>
          <span class="text-secondary">/mes</span>
        </div>
        <ul class="list-unstyled text-secondary">
          <li class="mb-1"><i class="bi bi-check-circle me-2" style="color:${colores[i]}"></i>Calidad ${p.calidad_maxima}</li>
          <li class="mb-1"><i class="bi bi-check-circle me-2" style="color:${colores[i]}"></i>${p.max_pantallas} pantalla(s)</li>
        </ul>
        <button class="btn btn-sm mt-3" style="background:${colores[i]};color:white">
          Seleccionar
        </button>
      </div>
    </div>` + "`" + `).join('');
}

// ── USUARIOS ──────────────────────────────────────────────────────────────────
async function cargarUsuarios() {
  const r = await fetch(API+'/api/usuarios');
  const data = await r.json();
  const colPlan = {'Básico':'badge-basico','Estándar':'badge-estandar','Premium':'badge-premium'};
  document.getElementById('tablaUsuarios').innerHTML = (data.datos||[]).map(u => ` + "`" + `
    <tr>
      <td>${u.id}</td>
      <td><i class="bi bi-person-circle me-2 text-secondary"></i>${u.nombre}</td>
      <td class="text-secondary">${u.email}</td>
      <td><span class="badge ${colPlan[u.plan]||'bg-secondary'}">${u.plan}</span></td>
      <td><span class="badge ${u.activo?'bg-success':'bg-danger'}">${u.activo?'Activo':'Inactivo'}</span></td>
    </tr>` + "`" + `).join('');
}

function mostrarFormUsuario() {
  const f = document.getElementById('formUsuario');
  f.style.display = f.style.display==='none' ? 'block' : 'none';
}

async function registrarUsuario() {
  const body = { nombre: document.getElementById('uNombre').value,
                 email: document.getElementById('uEmail').value,
                 password: document.getElementById('uPassword').value,
                 plan: document.getElementById('uPlan').value };
  const r = await fetch(API+'/api/usuarios', { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(body) });
  const data = await r.json();
  mostrarAlerta(data.mensaje, data.ok ? 'success' : 'danger');
  if (data.ok) { document.getElementById('formUsuario').style.display='none'; cargarUsuarios(); }
}

// ── REPORTES ──────────────────────────────────────────────────────────────────
async function cargarReportes() {
  const r = await fetch(API+'/api/reportes');
  const data = await r.json();
  if (data.datos) {
    document.getElementById('reporteContenido').textContent = data.datos.contenido || '';
    document.getElementById('reporteUsuarios').textContent  = data.datos.usuarios  || '';
  }
}

// ── LOGIN ─────────────────────────────────────────────────────────────────────
async function iniciarSesion() {
  const body = { email: document.getElementById('loginEmail').value,
                 password: document.getElementById('loginPassword').value };
  const r = await fetch(API+'/api/login', { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify(body) });
  const data = await r.json();
  if (data.ok) {
    usuarioActual = data.datos;
    document.getElementById('userBadge').classList.remove('d-none');
    document.getElementById('userNameBadge').textContent = usuarioActual.nombre;
    document.getElementById('userPlanBadge').textContent = usuarioActual.plan;
    mostrarAlerta('¡Bienvenido, ' + usuarioActual.nombre + '!', 'success');
  } else {
    mostrarAlerta(data.mensaje, 'danger');
  }
}

function cerrarSesion() {
  usuarioActual = null;
  document.getElementById('userBadge').classList.add('d-none');
  mostrarAlerta('Sesión cerrada', 'info');
}

// Carga inicial
cargarCatalogo();
</script>
</body>
</html>
`

// RegistrarFrontend agrega la ruta "/" que sirve la interfaz web.
func RegistrarFrontend() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl := template.Must(template.New("index").Parse(paginaHTML))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, nil)
	})
}
