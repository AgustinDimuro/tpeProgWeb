package db

import (
	"database/sql"
	"time"
	db "tpeProgWeb/db/sqlc"
	"tpeProgWeb/domain"
)

type DBReservationRepository struct {
	db *sql.DB
}

func NewDBReservationRepositoryADM(db *sql.DB) *DBReservationRepository {
	return &DBReservationRepository{db: db}
}

func (repo *DBReservationRepository) CreateReservation(reservation *domain.Reservation) error {
	queries := db.New(repo.db)
	_, err := queries.CreateReservation(ctx,
		db.CreateReservationParams{
			CabinID: int32(reservation.CabinID),
			Fecha:   reservation.Fecha,
		},
	)
	return err
}

func (repo *DBReservationRepository) GetAllReservations() ([]*domain.Reservation, error) {
	queries := db.New(repo.db)
	dbReservations, err := queries.ListReservations(ctx)
	if err != nil {
		return nil, err
	}

	var reservations []*domain.Reservation
	for _, dbRes := range dbReservations {
		reservation := domain.Reservation{
			ID:      int64(dbRes.ID),
			CabinID: int64(dbRes.CabinID),
			Fecha:   dbRes.Fecha,
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, nil
}
func (repo *DBReservationRepository) GetAllReservationsByCabinID(CabinID int64) ([]*domain.Reservation, error) {
	queries := db.New(repo.db)
	dbReservations, err := queries.ListReservationsByCabin(ctx, int32(CabinID))
	if err != nil {
		return nil, err
	}

	var reservations []*domain.Reservation
	for _, dbRes := range dbReservations {
		reservation := domain.Reservation{
			ID:      int64(dbRes.ID),
			CabinID: int64(dbRes.CabinID),
			Fecha:   dbRes.Fecha,
		}
		reservations = append(reservations, &reservation)
	}
	return reservations, nil
}

func (repo *DBReservationRepository) DeleteReservationByDate(CabinID int64, fecha time.Time) error {
	queries := db.New(repo.db)
	return queries.DeleteReservation(ctx, db.DeleteReservationParams{
		CabinID: int32(CabinID),
		Fecha:   fecha,
	})
}

func (repo *DBReservationRepository) ChangeReservationDate(cabinID int64, oldDate time.Time, newDate time.Time) error {
	queries := db.New(repo.db)
	_, err := queries.UpdateReservation(ctx, db.UpdateReservationParams{
		CabinID:  int32(cabinID),
		NewFecha: newDate,
		Fecha:    oldDate,
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBReservationRepository) GetReservationByDate(fecha time.Time) (*domain.Reservation, error) {
	queries := db.New(repo.db)

	dbRes, err := queries.GetReservationByFecha(ctx, fecha)
	if err != nil {
		return nil, err
	}

	reservation := &domain.Reservation{
		ID:      int64(dbRes.ID),
		CabinID: int64(dbRes.CabinID),
		Fecha:   dbRes.Fecha,
	}

	return reservation, nil
}
