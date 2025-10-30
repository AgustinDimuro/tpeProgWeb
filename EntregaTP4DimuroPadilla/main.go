package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	db "tpeProgWeb/db/sqlc"

	_ "github.com/lib/pq"
)

const dbSource = "postgres://user:password@localhost:5432/mydb?sslmode=disable"

var ctx context.Context
var dbQueries *db.Queries

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	cabin_id := r.URL.Query().Get("cabin_id")
	fecha := r.URL.Query().Get("fecha")
	fechaNueva := r.URL.Query().Get("fecha_nueva")

	switch r.Method {
	case http.MethodGet:
		getReservation(w, r, fecha)
	case http.MethodPost:
		createReservation(w, r, cabin_id, fecha)
	case http.MethodPut:
		updateReservation(w, r, cabin_id, fecha, fechaNueva)
	case http.MethodDelete:
		deleteReservation(w, r, cabin_id, fecha)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func reservationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllReservations(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func getAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := dbQueries.ListReservations(ctx)
	if err != nil {
		http.Error(w, "Error al obtener las reservas", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func deleteReservation(w http.ResponseWriter, r *http.Request, cabin_id, fecha string) {
	fechaParsed, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		http.Error(w, "La fecha no cumple el formato adecuado", http.StatusNotFound)
		return
	}
	cabinID, err := strconv.Atoi(cabin_id)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero", http.StatusBadRequest)
		return
	}
	err = dbQueries.DeleteReservation(ctx, db.DeleteReservationParams{
		CabinID: int32(cabinID),
		Fecha:   fechaParsed,
	})
	if err != nil {
		http.Error(w, "Error al eliminar la reserva", http.StatusInternalServerError)
		return
	}
}

func updateReservation(w http.ResponseWriter, r *http.Request, cabin_id, fecha, fechaNueva string) {
	fechaParsed, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		http.Error(w, "La fecha no cumple el formato adecuado", http.StatusNotFound)
		return
	}
	fechaNuevaParsed, err := time.Parse("2006-01-02", fechaNueva)
	if err != nil {
		http.Error(w, "La nueva fecha no cumple el formato adecuado", http.StatusNotFound)
		return
	}

	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)

	if fechaNuevaParsed.Before(todayUTC) {
		http.Error(w, "La nueva fecha de la reserva debe ser de hoy o a futuro", http.StatusBadRequest)
		return
	}

	cabinID, err := strconv.Atoi(cabin_id)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero", http.StatusBadRequest)
		return
	}
	reservation, err := dbQueries.UpdateReservation(ctx, db.UpdateReservationParams{
		CabinID:  int32(cabinID),
		Fecha:    fechaParsed,
		NewFecha: fechaNuevaParsed,
	})
	if err != nil {
		http.Error(w, "Error al actualizar la reserva", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func createReservation(w http.ResponseWriter, r *http.Request, cabin_id string, fecha string) {
	fechaParsed, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		http.Error(w, "La fecha no cumple el formato adecuado", http.StatusNotFound)
		return
	}

	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)

	if fechaParsed.Before(todayUTC) {
		http.Error(w, "La fecha de la reserva debe ser de hoy o a futuro", http.StatusBadRequest)
		return
	}

	cabinID, err := strconv.Atoi(cabin_id)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero", http.StatusBadRequest)
		return
	}
	reservation, err := dbQueries.CreateReservation(ctx, db.CreateReservationParams{
		CabinID: int32(cabinID),
		Fecha:   fechaParsed,
	})
	if err != nil {
		http.Error(w, "Error al crear la reserva", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func getReservation(w http.ResponseWriter, r *http.Request, fecha string) {
	fechaPased, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		http.Error(w, "La fecha no cumple el formato adecuado", http.StatusNotFound)
		return
	}
	reservation, err := dbQueries.GetReservationByFecha(ctx, fechaPased)
	if err != nil {
		http.Error(w, "Error al obtener la reserva", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func cabinHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	emailContact := r.URL.Query().Get("emailContact")
	phoneContact := r.URL.Query().Get("phoneContact")
	password := r.URL.Query().Get("password")

	switch r.Method {
	case http.MethodGet:
		getCabin(w, r, id)
	case http.MethodPut:
		updateCabin(w, r, id, emailContact, phoneContact, password)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func updateCabin(w http.ResponseWriter, r *http.Request, id, emailContact, phoneContact, password string) {
	cabinID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero", http.StatusBadRequest)
		return
	}
	cabin, err := dbQueries.UpdateCabin(ctx, db.UpdateCabinParams{
		ID:           int32(cabinID),
		EmailContact: emailContact,
		PhoneContact: phoneContact,
		Password:     password,
	})
	if err != nil {
		http.Error(w, "Error al actualizar la cabaña", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cabin)
}

func getCabin(w http.ResponseWriter, r *http.Request, id string) {
	cabinID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero", http.StatusBadRequest)
		return
	}
	cabin, err := dbQueries.GetCabin(ctx, int32(cabinID))
	if err != nil {
		http.Error(w, "Error al obtener la cabaña", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cabin)
}
func main() {
	// Conectar a la base
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("No se pudo conectar a la base:", err)
	}

	defer conn.Close()

	dbQueries = db.New(conn)
	ctx = context.Background()

	// 1. Crear una Cabin
	cabin, err := dbQueries.CreateCabin(ctx, db.CreateCabinParams{
		EmailContact: "contacto@ejemplo.com",
		PhoneContact: "123456789",
		Password:     "secreta",
	})
	if err != nil {
		log.Fatal("Error creando cabin:", err)
	}
	//	fmt.Printf("Cabin creada: %+v\n", cabin)

	cabin2, err := dbQueries.CreateCabin(ctx, db.CreateCabinParams{
		EmailContact: "contacto@ejemplo.com",
		PhoneContact: "123456789",
		Password:     "secreta",
	})
	if err != nil {
		log.Fatal("Error creando cabin:", err)
	}
	//	fmt.Printf("Cabin creada: %+v\n", cabin2)

	// 2. Crear una Reservation
	res, err := dbQueries.CreateReservation(ctx, db.CreateReservationParams{
		CabinID: cabin.ID,
		Fecha:   time.Now().AddDate(0, 0, 8), // reserva dentro de 7 días
	})
	if err != nil {
		log.Fatal("Error creando reservation:", err)
	}

	res2, err := dbQueries.CreateReservation(ctx, db.CreateReservationParams{
		CabinID: cabin2.ID,
		Fecha:   time.Now().AddDate(0, 0, 17), // reserva dentro de 7 días
	})
	if err != nil {
		log.Fatal("Error creando reservation:", err)
	}

	fmt.Printf("Reservation creada: %+v\n", res)
	fmt.Printf("Reservation creada: %+v\n", res2)
	/*
		// 3. Listar todas las cabins
		cabins, err := dbQueries.ListCabins(ctx)
		if err != nil {
			log.Fatal("Error listando cabins:", err)
		}
		fmt.Println("Todas las cabins:")
		for _, c := range cabins {
			fmt.Printf(" - %+v\n", c)
		}

		// 4. Listar todas las reservations
		reservations, err := dbQueries.ListReservations(ctx)
		if err != nil {
			log.Fatal("Error listando reservations:", err)
		}
		fmt.Println("Todas las reservations:")
		for _, r := range reservations {
			fmt.Printf(" - %+v\n", r)
		}

		// 5. Probar disponibilidad de fecha
		fecha := time.Now().AddDate(0, 0, 7)
		disponible, err := dbQueries.IsFechaDisponible(ctx, fecha)
		if err != nil {
			log.Fatal("Error verificando disponibilidad:", err)
		}
		fmt.Printf("¿Fecha %s disponible?: %v\n", fecha.Format("2006-01-02"), disponible)
	*/
	http.HandleFunc("/reservation", reservationHandler)
	http.HandleFunc("/reservations", reservationsHandler)

	http.HandleFunc("/cabin", cabinHandler)

	// Sirve archivos estáticos (index.html, app.js, styles.css)
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: logica\n")

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
