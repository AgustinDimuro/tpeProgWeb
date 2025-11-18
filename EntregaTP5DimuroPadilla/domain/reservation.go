package domain

import "time"

type Reservation struct {
	ID      int64
	CabinID int64
	Fecha   time.Time
}

type ReservationADMRepository interface {
	CreateReservation(reservation *Reservation) error
	GetAllReservations() ([]*Reservation, error)
	GetAllReservationsByCabinID(CabinID int64) ([]*Reservation, error)
	DeleteReservationByDate(CabinID int64, fecha time.Time) error
	ChangeReservationDate(CabinID int64, oldDate time.Time, newDate time.Time) error
	GetReservationByDate(fecha time.Time) (*Reservation, error)
}
type ReservationUserRepository interface {
	CreateReservation(reservation *Reservation) error
	GetAllReservations() ([]*Reservation, error)
	GetAllReservationsByCabinID(CabinID int64) ([]*Reservation, error)
	DeleteReservationByDate(CabinID int64, fecha time.Time) error
	ChangeReservationDate(CabinID int64, oldDate time.Time, newDate time.Time) error
	GetReservationByDate(fecha time.Time) (*Reservation, error)
}
