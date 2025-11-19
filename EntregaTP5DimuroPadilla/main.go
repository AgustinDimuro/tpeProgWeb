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

	// --- INSERCIÓN DE DATOS DE PRUEBA (Opcional: Puedes comentarlo si ya existen) ---
	fmt.Println("Insertando datos de prueba...")
	dbQueries := dbsqlc.New(conn)
	ctx := context.Background()

	cabin1, err := dbQueries.CreateCabin(ctx, dbsqlc.CreateCabinParams{
		EmailContact: "contacto@cabania1.com",
		PhoneContact: "123456789",
		Password:     "secreta",
	})
	if err != nil {
		log.Printf("No se pudo crear cabin1 (quizás ya existe): %v", err)
	} else {
		fmt.Printf("Cabaña 1 creada con ID: %d\n", cabin1.ID)
		_, err = dbQueries.CreateReservation(ctx, dbsqlc.CreateReservationParams{
			CabinID: cabin1.ID,
			Fecha:   time.Now().AddDate(0, 0, 10),
		})
		if err != nil {
			log.Printf("No se pudo crear reserva para cabin1: %v", err)
		} else {
			fmt.Println("Reserva 1 creada.")
		}
	}
	// -------------------------------------------------------------------------------

	cabinRepo := db.NewDBCabinRepository(conn)
	reservationRepo := db.NewDBReservationRepositoryADM(conn)

	cabinServiceAdm := application.NewCabinServiceADM(cabinRepo)
	cabinServiceUser := application.NewCabinServiceUser(cabinRepo)

	reservationServiceAdm := application.NewReservationServicesADM(reservationRepo)
	reservationServiceUser := application.NewReservationServicesUser(reservationRepo)

	adminHandler := ui.NewAdminHandler(reservationServiceAdm, cabinServiceAdm)
	userHandler := ui.NewUserHandler(reservationServiceUser, cabinServiceUser)

	// --- NUEVO: Inicializamos el AuthHandler ---
	authHandler := ui.NewAuthHandler(cabinServiceUser)

	// -----------------------------------------------------------------
	// RUTAS PÚBLICAS (Login y Logout)
	// -----------------------------------------------------------------
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.HandleLoginShow(w, r)
		} else if r.Method == http.MethodPost {
			authHandler.HandleLoginProcess(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", authHandler.HandleLogout)

	// -----------------------------------------------------------------
	// RUTAS PROTEGIDAS (Usuario) - Envueltas en AuthMiddleware
	// -----------------------------------------------------------------

	// Página Principal (Dashboard)
	http.HandleFunc("/", ui.AuthMiddleware(userHandler.HandleShowMainPage))

	// Calendario
	http.HandleFunc("/calendario", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.HandleShowCalendar(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}))

	// Gestión de Reservas
	http.HandleFunc("/reservations", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Get("fecha") != "" {
				userHandler.GetReservationByDateHandler(w, r)
			} else {
				userHandler.GetReservationByDateHandler(w, r)
			}
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
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/reservations/delete", ui.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.DeleteReservationHandler(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}))

	// -----------------------------------------------------------------
	// RUTAS DE ADMINISTRACIÓN (Admin)
	// Nota: Por ahora están sin protección o podrías usar el mismo middleware si aplica
	// -----------------------------------------------------------------
	http.HandleFunc("/admin/cabins", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			adminHandler.UpdateCabinHandler(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/admin/reservations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			adminHandler.GetAllReservationsHandler(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
