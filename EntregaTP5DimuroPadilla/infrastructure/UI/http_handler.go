package ui

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"tpeProgWeb/application"
	"tpeProgWeb/domain"
)

var ErrInvalidReservationDate = errors.New("la fecha de la reserva debe ser de hoy o a futuro")

type AdminHandler struct {
	reservationServiceADM *application.ReservationServicesADM
	cabinServiceADM       *application.CabinServicesADM
}
type UserHandler struct {
	reservationServiceUser *application.ReservationServicesUser
	cabinServiceUser       *application.CabinServicesUser
}

func NewAdminHandler(reservationServiceADM *application.ReservationServicesADM, cabinServiceADM *application.CabinServicesADM) *AdminHandler {
	return &AdminHandler{
		reservationServiceADM: reservationServiceADM,
		cabinServiceADM:       cabinServiceADM,
	}
}

func NewUserHandler(reservationServiceUser *application.ReservationServicesUser, cabinServiceUser *application.CabinServicesUser) *UserHandler {
	return &UserHandler{
		reservationServiceUser: reservationServiceUser,
		cabinServiceUser:       cabinServiceUser,
	}
}

func (h *AdminHandler) GetCabinByIDHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Se requiere el parámetro 'id'", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	cabin, err := h.cabinServiceADM.GetCabinByID(id)
	if err != nil {

		http.Error(w, "Error al obtener la cabaña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cabin)
}

func (h *AdminHandler) UpdateCabinHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Se requiere el parámetro 'id'", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "El ID de la cabaña debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	var cabinData domain.Cabin
	if err := json.NewDecoder(r.Body).Decode(&cabinData); err != nil {
		http.Error(w, "Cuerpo de la petición inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	cabinData.ID = id

	if err := h.cabinServiceADM.UpdateCabin(&cabinData); err != nil {
		http.Error(w, "Error al actualizar la cabaña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cabaña actualizada con éxito"})
}

func (h *AdminHandler) GetAllReservationsHandler(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.reservationServiceADM.GetAllReservations()
	if err != nil {
		http.Error(w, "Error al obtener las reservas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}

func (h *UserHandler) CreateReservationHandler(w http.ResponseWriter, r *http.Request) {

	var newRes domain.Reservation
	if err := json.NewDecoder(r.Body).Decode(&newRes); err != nil {
		http.Error(w, "Cuerpo de la petición inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.reservationServiceUser.CreateReservation(&newRes); err != nil {

		if errors.Is(err, application.ErrInvalidReservationDate) {

			http.Error(w, err.Error(), http.StatusBadRequest)

		} else {
			http.Error(w, "Error interno al crear la reserva: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reserva creada con éxito"})
}

func (h *UserHandler) UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {

	cabinIDStr := r.URL.Query().Get("cabin_id")
	oldDateStr := r.URL.Query().Get("fecha")

	if cabinIDStr == "" || oldDateStr == "" {
		http.Error(w, "Se requieren los parámetros 'cabin_id' y 'fecha'", http.StatusBadRequest)
		return
	}

	cabinID, err := strconv.ParseInt(cabinIDStr, 10, 64)
	if err != nil {
		http.Error(w, "El 'cabin_id' debe ser un número entero", http.StatusBadRequest)
		return
	}

	oldDate, err := time.Parse("2006-01-02", oldDateStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha' debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	var updateReq domain.Reservation
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Cuerpo de la petición inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	newDate := updateReq.Fecha

	if err := h.reservationServiceUser.ChangeReservationDate(cabinID, oldDate, newDate); err != nil {

		if errors.Is(err, application.ErrInvalidReservationDate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Error interno al actualizar la reserva: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reserva actualizada con éxito"})
}

func (h *UserHandler) DeleteReservationHandler(w http.ResponseWriter, r *http.Request) {

	cabinIDStr := r.URL.Query().Get("cabin_id")
	dateStr := r.URL.Query().Get("fecha")

	if cabinIDStr == "" || dateStr == "" {
		http.Error(w, "Se requieren los parámetros 'cabin_id' y 'fecha'", http.StatusBadRequest)
		return
	}

	cabinID, err := strconv.ParseInt(cabinIDStr, 10, 64)
	if err != nil {
		http.Error(w, "El 'cabin_id' debe ser un número entero", http.StatusBadRequest)
		return
	}

	fecha, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha' debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	if err := h.reservationServiceUser.DeleteReservationByDate(cabinID, fecha); err != nil {
		http.Error(w, "Error al eliminar la reserva: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reserva eliminada con éxito"})
}

func (h *UserHandler) GetReservationByDateHandler(w http.ResponseWriter, r *http.Request) {

	dateStr := r.URL.Query().Get("fecha")
	if dateStr == "" {
		http.Error(w, "Se requiere el parámetro 'fecha'", http.StatusBadRequest)
		return
	}

	fecha, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha' debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	reservation, err := h.reservationServiceUser.GetReservationByDate(fecha)
	if err != nil {
		http.Error(w, "Reserva no encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservation)
}
