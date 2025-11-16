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

	if err = conn.Ping(); err != nil {
		log.Fatal("No se pudo hacer ping a la base de datos:", err)
	}

	fmt.Println("Conexión a la base de datos exitosa.")

	fmt.Println("Insertando datos de prueba...")

	// 1. Preparamos las consultas de sqlc
	dbQueries := dbsqlc.New(conn)
	ctx := context.Background()

	// 2. Creamos una Cabaña 1
	cabin1, err := dbQueries.CreateCabin(ctx, dbsqlc.CreateCabinParams{
		EmailContact: "contacto@cabania1.com",
		PhoneContact: "123456789",
		Password:     "secreta",
	})
	if err != nil {
		// Usamos Printf en lugar de Fatal. Si la cabaña ya existe,
		// la app no se detendrá, solo lo informará.
		log.Printf("No se pudo crear cabin1 (quizás ya existe): %v", err)
	} else {
		fmt.Printf("Cabaña 1 creada con ID: %d\n", cabin1.ID)

		// 3. Creamos una Reserva para la Cabaña 1
		// (Solo si la cabaña se creó recién)
		_, err = dbQueries.CreateReservation(ctx, dbsqlc.CreateReservationParams{
			CabinID: cabin1.ID,
			Fecha:   time.Now().AddDate(0, 0, 10), // Reserva para 10 días en el futuro
		})
		if err != nil {
			log.Printf("No se pudo crear reserva para cabin1: %v", err)
		} else {
			fmt.Println("Reserva 1 creada.")
		}
	}

	// 4. Creamos una Cabaña 2
	cabin2, err := dbQueries.CreateCabin(ctx, dbsqlc.CreateCabinParams{
		EmailContact: "info@lagoazul.com",
		PhoneContact: "987654321",
		Password:     "clave123",
	})
	if err != nil {
		log.Printf("No se pudo crear cabin2 (quizás ya existe): %v", err)
	} else {
		fmt.Printf("Cabaña 2 creada con ID: %d\n", cabin2.ID)
	}

	fmt.Println("Datos de prueba insertados.")

	cabinRepo := db.NewDBCabinRepository(conn)
	reservationRepo := db.NewDBReservationRepositoryADM(conn)

	cabinServiceAdm := application.NewCabinServiceADM(cabinRepo)
	cabinServiceUser := application.NewCabinServiceUser(cabinRepo)

	reservationServiceAdm := application.NewReservationServicesADM(reservationRepo)
	reservationServiceUser := application.NewReservationServicesUser(reservationRepo)

	adminHandler := ui.NewAdminHandler(reservationServiceAdm, cabinServiceAdm)
	userHandler := ui.NewUserHandler(reservationServiceUser, cabinServiceUser)

	// -----------------------------------------------------------------
	// PASO 5: CONFIGURAR EL ROUTER (UI)
	// -----------------------------------------------------------------
	// main.go ya no tiene lógica de negocio, solo enruta las peticiones
	// a los handlers correctos.

	// ---- Rutas de Cabañas (Admin) ----
	http.HandleFunc("/admin/cabins", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adminHandler.GetCabinByIDHandler(w, r)
		case http.MethodPut:
			adminHandler.UpdateCabinHandler(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// ---- Rutas de Reservas (Admin) ----
	http.HandleFunc("/admin/reservations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adminHandler.GetAllReservationsHandler(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// ---- Rutas de Reservas (User) ----
	// Esta ruta maneja todo lo relacionado con las reservas de usuario
	http.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Verificamos si nos piden una reserva por fecha
			if r.URL.Query().Get("fecha") != "" {
				userHandler.GetReservationByDateHandler(w, r)
			} else {
				// Si no, podríamos manejar "Get por ID de cabaña" aquí
				// (Esa función ya la tienes en tu servicio: GetAllReservationsByCabinID)
				userHandler.GetReservationByDateHandler(w, r)
				//http.Error(w, "Parámetro 'fecha' (para GET) o método no implementado", http.StatusBadRequest)
			}
		case http.MethodPost:
			userHandler.CreateReservationHandler(w, r)
		case http.MethodPut:
			userHandler.UpdateReservationHandler(w, r)
		case http.MethodDelete:
			userHandler.DeleteReservationHandler(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// -----------------------------------------------------------------
	// PASO 6: SERVIR ARCHIVOS ESTÁTICOS E INICIAR EL SERVIDOR
	// -----------------------------------------------------------------

	// Sirve archivos estáticos (index.html, app.js, styles.css) desde el directorio raíz
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
