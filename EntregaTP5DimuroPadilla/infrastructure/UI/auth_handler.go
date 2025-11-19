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

// Clave secreta para firmar los tokens (En prod debe ir en variable de entorno)
var jwtKey = []byte("mi_clave_super_secreta_del_quincho")

// Defines los "Claims" (datos) que irán dentro del token
type Claims struct {
	CabinID int64 `json:"cabin_id"`
	jwt.RegisteredClaims
}

type AuthHandler struct {
	CabinService *application.CabinServicesUser
}

func NewAuthHandler(service *application.CabinServicesUser) *AuthHandler {
	return &AuthHandler{CabinService: service}
}

// 1. Mostrar Formulario de Login
func (h *AuthHandler) HandleLoginShow(w http.ResponseWriter, r *http.Request) {
	views.Login().Render(r.Context(), w)
}

// 2. Procesar Login (POST)
func (h *AuthHandler) HandleLoginProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error form", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("cabin_id")
	password := r.FormValue("password")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Llamamos al servicio de aplicación
	cabin, err := h.CabinService.Authenticate(id, password)
	if err != nil {
		// Login fallido
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	} // Paso 2: Aplicar regla de negocio (comparación)
	// Aquí en el futuro cambiaremos "==" por bcrypt.CompareHashAndPassword

	// --- Generar JWT ---
	expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 1 día
	claims := &Claims{
		CabinID: cabin.ID,
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

	// --- Guardar en Cookie HttpOnly ---
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true, // Importante: No accesible por JS
		Path:     "/",  // Disponible en toda la app
		SameSite: http.SameSiteStrictMode,
	})

	// Redirigir al Home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 3. Logout
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Sobrescribimos la cookie con fecha expirada
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// 4. Middleware de Protección
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtenemos la cookie
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// No hay cookie, redirigir a login
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Error de autenticación", http.StatusBadRequest)
			return
		}

		tokenStr := c.Value
		claims := &Claims{}

		// Parsear y validar token
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Token válido: Pasamos al siguiente handler
		// Opcional: Podríamos inyectar el ID en el contexto si lo necesitamos luego
		ctx := context.WithValue(r.Context(), "cabinID", claims.CabinID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
