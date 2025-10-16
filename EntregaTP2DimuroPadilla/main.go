package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	db "TPEProgWebEntregaTP2/db/sqlc"

	_ "github.com/lib/pq"
)

const dbSource = "postgres://user:password@localhost:5432/mydb?sslmode=disable"

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
		Fecha:   time.Now().AddDate(0, 0, 7), // reserva dentro de 7 días
	})
	if err != nil {
		log.Fatal("Error creando reservation:", err)
	}
	fmt.Printf("Reservation creada: %+v\n", res)

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
}
