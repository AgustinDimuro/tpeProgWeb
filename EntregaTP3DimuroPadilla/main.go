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
var dbQueries db.Queries

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
	reservation, err := dbQueries.GetReservationByFecha(ctx, fechaPased)
	if err != nil {
		http.Error(w, "La fecha no cumple el formato adecuado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func cabinHandler(w http.ResponseWriter, r *http.Request) {

}
func main() {
	// Conectar a la base
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("No se pudo conectar a la base:", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	ctx := context.Background()

	// 1. Crear una Cabin
	cabin, err := queries.CreateCabin(ctx, db.CreateCabinParams{
		EmailContact: "contacto@ejemplo.com",
		PhoneContact: "123456789",
		Password:     "secreta",
	})
	if err != nil {
		log.Fatal("Error creando cabin:", err)
	}
	fmt.Printf("Cabin creada: %+v\n", cabin)

	// 2. Crear una Reservation
	res, err := queries.CreateReservation(ctx, db.CreateReservationParams{
		CabinID: cabin.ID,
		Fecha:   time.Now().AddDate(0, 0, 8), // reserva dentro de 7 días
	})
	if err != nil {
		log.Fatal("Error creando reservation:", err)
	}
	fmt.Printf("Reservation creada: %+v\n", res)
	/*
		// 3. Listar todas las cabins
		cabins, err := queries.ListCabins(ctx)
		if err != nil {
			log.Fatal("Error listando cabins:", err)
		}
		fmt.Println("Todas las cabins:")
		for _, c := range cabins {
			fmt.Printf(" - %+v\n", c)
		}

		// 4. Listar todas las reservations
		reservations, err := queries.ListReservations(ctx)
		if err != nil {
			log.Fatal("Error listando reservations:", err)
		}
		fmt.Println("Todas las reservations:")
		for _, r := range reservations {
			fmt.Printf(" - %+v\n", r)
		}

		// 5. Probar disponibilidad de fecha
		fecha := time.Now().AddDate(0, 0, 7)
		disponible, err := queries.IsFechaDisponible(ctx, fecha)
		if err != nil {
			log.Fatal("Error verificando disponibilidad:", err)
		}
		fmt.Printf("¿Fecha %s disponible?: %v\n", fecha.Format("2006-01-02"), disponible)
	*/
	http.HandleFunc("/reservations", reservationHandler)
	http.HandleFunc("/cabins", cabinHandler)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: logica\n")

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
