package ui

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"tpeProgWeb/application"
	"tpeProgWeb/views"

	"github.com/golang-jwt/jwt/v5"
)

// CLAVE SECRETA (En producción esto debe ir en variables de entorno)
var jwtKey = []byte("mi_clave_super_secreta_del_quincho")

// Claims define qué datos guardamos DENTRO del token
type Claims struct {
	CabinID int64  `json:"cabin_id"`
	Role    string `json:"role"` // <--- Importante: Aquí viaja el rol
	jwt.RegisteredClaims
}

type AuthHandler struct {
	CabinService *application.CabinServicesUser
}

func NewAuthHandler(service *application.CabinServicesUser) *AuthHandler {
	return &AuthHandler{CabinService: service}
}

// 1. GET /login - Muestra el formulario
func (h *AuthHandler) HandleLoginShow(w http.ResponseWriter, r *http.Request) {
	// Si ya tiene cookie válida, lo mandamos directo al home
	if _, err := r.Cookie("token"); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	views.Login().Render(r.Context(), w)
}

// 2. POST /login - Procesa las credenciales
func (h *AuthHandler) HandleLoginProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar formulario", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("cabin_id")
	password := r.FormValue("password")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		views.Login().Render(r.Context(), w) // Podrías pasar un mensaje de error aquí
		return
	}

	// --- LLAMADA A LA LÓGICA DE NEGOCIO ---
	// Esto verifica ID, Pass y nos devuelve el ROL
	cabin, err := h.CabinService.Authenticate(id, password)
	if err != nil {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// --- GENERACIÓN DEL JWT ---
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		CabinID: cabin.ID,
		Role:    cabin.Role, // <--- Guardamos el Rol en el Token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error interno firmando token", http.StatusInternalServerError)
		return
	}

	// --- SETEAR COOKIE HTTPONLY ---
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true, // Seguridad: JS no puede leer esto
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	// --- REDIRECCIÓN SEGÚN ROL ---
	if cabin.Role == "admin" {
		http.Redirect(w, r, "/admin/reservations", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// 3. GET /logout - Borra la cookie
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Fecha en el pasado
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// --- MIDDLEWARES ---

// Middleware GENERAL (Solo verifica que estés logueado)
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validar Cookie
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Error de autenticación", http.StatusBadRequest)
			return
		}

		// Validar Token
		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Inyectar datos en el contexto por si los necesitamos en el handler
		ctx := context.WithValue(r.Context(), "userRole", claims.Role)
		ctx = context.WithValue(ctx, "cabinID", claims.CabinID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Middleware ADMIN (Verifica logueo Y Rol Admin)
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// --- VERIFICACIÓN DE ROL ---
		if claims.Role != "admin" {
			http.Error(w, "Acceso Prohibido: Requiere permisos de Administrador", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "userRole", claims.Role)
		ctx = context.WithValue(ctx, "cabinID", claims.CabinID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
