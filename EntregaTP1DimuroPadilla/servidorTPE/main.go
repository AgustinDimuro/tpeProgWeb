package main

import (
	"fmt"
	"net/http"
)

func main() {
	staticDir := "./static"

	// Servidor de archivos estáticos (index.html, reservar.html, etc.)
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	// Handler de login
	http.HandleFunc("/login", handleLogin)

	// Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// Inicia el servidor
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}

// Handler para procesar login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer datos del formulario
	cabinID := r.FormValue("id")
	password := r.FormValue("password")

	// Ejemplo simplificado: hardcodeamos un usuario válido
	if cabinID == "1" && password == "1234" {
		// Login correcto → redirigir a reservas
		http.Redirect(w, r, "/reservar.html", http.StatusSeeOther)
		return
	}

	// Si no coincide → error
	http.Error(w, "Número de cabaña o contraseña incorrectos", http.StatusUnauthorized)
}
