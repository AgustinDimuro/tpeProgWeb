package application

import (
	"time"
	"tpeProgWeb/domain"
)

type ReservationServicesADM struct {
	userRepository domain.ReservationADMRepository
}

func NewReservationServicesADM(userRepo domain.ReservationADMRepository) *ReservationServicesADM {
	return &ReservationServicesADM{
		userRepository: userRepo,
	}
}

func (service *ReservationServicesADM) CreateReservation(reservation *domain.Reservation) error {
	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)

	if reservation.Fecha.Before(todayUTC) {
		return ErrInvalidReservationDate
	}
	return service.userRepository.CreateReservation(reservation)
}

func (service *ReservationServicesADM) GetAllReservations() ([]*domain.Reservation, error) {
	return service.userRepository.GetAllReservations()
}

func (service *ReservationServicesADM) GetAllReservationsByCabinID(cabinID int64) ([]*domain.Reservation, error) {
	return service.userRepository.GetAllReservationsByCabinID(cabinID)
}
func (service *ReservationServicesADM) DeleteReservationByDate(CabinID int64, fecha time.Time) error {
	return service.userRepository.DeleteReservationByDate(CabinID, fecha)
}

func (service *ReservationServicesADM) ChangeReservationDate(CabinID int64, oldDate time.Time, newDate time.Time) error {
	nowUTC := time.Now().UTC()
	todayUTC := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)

	if newDate.Before(todayUTC) {
		return ErrInvalidReservationDate
	}
	return service.userRepository.ChangeReservationDate(CabinID, oldDate, newDate)
}

func (service *ReservationServicesADM) GetReservationByDate(fecha time.Time) (*domain.Reservation, error) {
	return service.userRepository.GetReservationByDate(fecha)
}
