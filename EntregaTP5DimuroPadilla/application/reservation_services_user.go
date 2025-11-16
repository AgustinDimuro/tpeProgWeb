package application

import (
	"errors"
	"time"
	"tpeProgWeb/domain"
)

var ErrInvalidReservationDate = errors.New("la fecha de la reserva debe ser de hoy o a futuro")

type ReservationServicesUser struct {
	userRepository domain.ReservationUserRepository
}

func NewReservationServicesUser(userRepo domain.ReservationUserRepository) *ReservationServicesUser {
	return &ReservationServicesUser{
		userRepository: userRepo,
	}
}

func (service *ReservationServicesUser) CreateReservation(reservation *domain.Reservation) error {
	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)

	if reservation.Fecha.Before(todayUTC) {
		return ErrInvalidReservationDate
	}
	return service.userRepository.CreateReservation(reservation)
}

func (service *ReservationServicesUser) GetAllReservationsByCabinID(cabinID int64) ([]*domain.Reservation, error) {
	return service.userRepository.GetAllReservationsByCabinID(cabinID)
}

func (service *ReservationServicesUser) DeleteReservationByDate(CabinID int64, fecha time.Time) error {
	return service.userRepository.DeleteReservationByDate(CabinID, fecha)
}

func (service *ReservationServicesUser) ChangeReservationDate(CabinID int64, oldDate time.Time, newDate time.Time) error {
	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)

	if newDate.Before(todayUTC) {
		return ErrInvalidReservationDate
	}
	return service.userRepository.ChangeReservationDate(CabinID, oldDate, newDate)
}

func (service *ReservationServicesUser) GetReservationByDate(fecha time.Time) (*domain.Reservation, error) {
	return service.userRepository.GetReservationByDate(fecha)
}
