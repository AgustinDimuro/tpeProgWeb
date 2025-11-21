package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"tpeProgWeb/application"
	dbsqlc "tpeProgWeb/db/sqlc"
	db "tpeProgWeb/infrastructure/DB"
	ui "tpeProgWeb/infrastructure/UI"

	_ "github.com/lib/pq"
)

const dbSource = "postgres://user:password@localhost:5432/mydb?sslmode=disable"

func main() {

	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("No se pudo conectar a la base:", err)
	}
	defer conn.Close()

	// Reintentos de conexión (útil para docker-compose)
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = conn.Ping()
		if err == nil {
			break
		}
		fmt.Printf("Intento %d/%d fallido: %v. Reintentando en 2s...\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err = conn.Ping(); err != nil {
		log.Fatal("No se pudo hacer ping a la base de datos:", err)
	}

	fmt.Println("Conexión a la base de datos exitosa.")

	// =================================================================
	// 1. DATOS DE PRUEBA (Actualizados con ROL)
	// =================================================================
	fmt.Println("Insertando datos de prueba...")
	dbQueries := dbsqlc.New(conn)
	ctx := context.Background()

	// Cabaña Usuario (Role: user)
	cabinUser, err := dbQueries.CreateCabin(ctx, dbsqlc.CreateCabinParams{
		EmailContact: "usuario@cabania.com",
		PhoneContact: "11111111",
		Password:     "1234", // Contraseña simple para probar
		Role:         "user", // <--- Rol Usuario
	})
	if err != nil {
		log.Printf("Nota: Cabaña Usuario ya existe o error: %v", err)
	} else {
		fmt.Printf(">> Creada Cabaña USUARIO (ID: %d, Pass: 1234)\n", cabinUser.ID)
	}

	// Cabaña Admin (Role: admin)
	cabinAdmin, err := dbQueries.CreateCabin(ctx, dbsqlc.CreateCabinParams{
		EmailContact: "admin@sistema.com",
		PhoneContact: "99999999",
		Password:     "admin", // Contraseña admin
		Role:         "admin", // <--- Rol Admin
	})
	if err != nil {
		log.Printf("Nota: Cabaña Admin ya existe o error: %v", err)
	} else {
		fmt.Printf(">> Creada Cabaña ADMIN (ID: %d, Pass: admin)\n", cabinAdmin.ID)
	}
	fmt.Println("Datos de prueba listos.")
	// =================================================================

	// 2. INICIALIZACIÓN DE CAPAS
	cabinRepo := db.NewDBCabinRepository(conn)
	reservationRepo := db.NewDBReservationRepositoryADM(conn)

	cabinServiceAdm := application.NewCabinServiceADM(cabinRepo)
	cabinServiceUser := application.NewCabinServiceUser(cabinRepo)

	reservationServiceAdm := application.NewReservationServicesADM(reservationRepo)
	reservationServiceUser := application.NewReservationServicesUser(reservationRepo)

	// Handlers
	adminHandler := ui.NewAdminHandler(reservationServiceAdm, cabinServiceAdm)
	userHandler := ui.NewUserHandler(reservationServiceUser, cabinServiceUser)
	authHandler := ui.NewAuthHandler(cabinServiceUser) // <--- Handler de Auth

	// 3. DEFINICIÓN DE RUTAS
	// =================================================================

	// --- RUTAS PÚBLICAS (Sin protección) ---
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.HandleLoginShow(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.HandleLoginProcess(w, r)
		}
	})

	http.HandleFunc("/logout", authHandler.HandleLogout)

	// --- RUTAS DE USUARIO (Protegidas con AuthMiddleware) ---
	// Cualquiera logueado (user o admin) puede ver esto, o podrías restringirlo solo a 'user' si quisieras.
	// Por ahora usamos AuthMiddleware genérico.

	http.HandleFunc("/", ui.AuthMiddleware(userHandler.HandleShowMainPage))

	http.HandleFunc("/calendario", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.HandleShowCalendar(w, r)
		}
	}))

	http.HandleFunc("/reservations", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetReservationByDateHandler(w, r)
		case http.MethodPost:
			userHandler.CreateReservationHandler(w, r)
		case http.MethodDelete:
			userHandler.DeleteReservationHandler(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/reservations/update", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.UpdateReservationHandler(w, r)
		}
	}))

	http.HandleFunc("/reservations/delete", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.DeleteReservationHandler(w, r)
		}
	}))

	// --- RUTAS DE ADMIN (Protegidas con AdminMiddleware) ---
	// Solo accesible si role == "admin"

	http.HandleFunc("/admin/reservations", ui.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			adminHandler.GetAllReservationsHandler(w, r)
		}
	}))

	http.HandleFunc("/admin/cabins", ui.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			adminHandler.UpdateCabinHandler(w, r)
		}
	}))

	// =================================================================
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
