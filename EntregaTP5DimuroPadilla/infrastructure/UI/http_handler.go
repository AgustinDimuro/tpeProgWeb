package ui

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"tpeProgWeb/application"
	"tpeProgWeb/domain"
	"tpeProgWeb/views"
)

var ErrInvalidReservationDate = errors.New("la fecha de la reserva debe ser de hoy o a futuro")

type AdminHandler struct {
	ReservationServiceADM *application.ReservationServicesADM
	CabinServiceADM       *application.CabinServicesADM
}
type UserHandler struct {
	ReservationServiceUser *application.ReservationServicesUser
	CabinServiceUser       *application.CabinServicesUser
}

func (h *UserHandler) HandleShowMainPage(w http.ResponseWriter, r *http.Request) {
	// Obtenemos todas las reservaciones para el usuario y las pasamos al Layout
	// Aqui habria que definir de qué cabaña son las reservas a mostrar
	cabinID := int64(1) // Por ejemplo, la cabaña con ID 1
	reservations, err := h.ReservationServiceUser.GetAllReservationsByCabinID(cabinID)
	if err != nil {
		log.Printf("Error al obtener reservaciones: %v", err)
		http.Error(w, "No se pudieron cargar las reservaciones", http.StatusInternalServerError)
		return
	}

	component := views.Layout(reservations)
	component.Render(r.Context(), w)
}

func NewAdminHandler(reservationServiceADM *application.ReservationServicesADM, cabinServiceADM *application.CabinServicesADM) *AdminHandler {
	return &AdminHandler{
		ReservationServiceADM: reservationServiceADM,
		CabinServiceADM:       cabinServiceADM,
	}
}

func NewUserHandler(reservationServiceUser *application.ReservationServicesUser, cabinServiceUser *application.CabinServicesUser) *UserHandler {
	return &UserHandler{
		ReservationServiceUser: reservationServiceUser,
		CabinServiceUser:       cabinServiceUser,
	}
}

func (h *UserHandler) CreateReservationHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	cabinIDStr := r.FormValue("cabin_id")
	fechaStr := r.FormValue("fecha")

	cabinID, err := strconv.ParseInt(cabinIDStr, 10, 64)
	if err != nil {
		http.Error(w, "El 'cabin_id' debe ser un número", http.StatusBadRequest)
		return
	}
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha' debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	var newRes = domain.Reservation{
		CabinID: cabinID,
		Fecha:   fecha,
	}

	if err := h.ReservationServiceUser.CreateReservation(&newRes); err != nil {

		if errors.Is(err, application.ErrInvalidReservationDate) {

			http.Error(w, err.Error(), http.StatusBadRequest)

		} else {
			http.Error(w, "Error interno al crear la reserva: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (h *UserHandler) UpdateReservationHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear el formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	cabinIDStr := r.FormValue("cabin_id")
	oldDateStr := r.FormValue("fecha")
	newDateStr := r.FormValue("fecha_nueva")

	cabinID, err := strconv.ParseInt(cabinIDStr, 10, 64)
	if err != nil {
		http.Error(w, "El 'cabin_id' debe ser un número", http.StatusBadRequest)
		return
	}
	oldDate, err := time.Parse("2006-01-02", oldDateStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha' (actual) debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}
	newDate, err := time.Parse("2006-01-02", newDateStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha_nueva' debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	if err := h.ReservationServiceUser.ChangeReservationDate(cabinID, oldDate, newDate); err != nil {

		if errors.Is(err, application.ErrInvalidReservationDate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Error interno al actualizar la reserva: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
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

	cabin, err := h.CabinServiceADM.GetCabinByID(id)
	if err != nil {

		http.Error(w, "Error al obtener la cabaña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cabin)
}
*/

func (h *AdminHandler) GetCabinByIDHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el ID de la URL (ej: /admin/cabins/edit?id=5)
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

	// 2. Buscar los datos actuales en la BD
	cabin, err := h.CabinServiceADM.GetCabinByID(id)
	if err != nil {
		http.Error(w, "Error al obtener la cabaña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Renderizar la vista de edición (templ) en lugar de JSON
	component := views.AdminEditCabin(cabin)
	component.Render(r.Context(), w)
}

/*
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

	if err := h.CabinServiceADM.UpdateCabin(&cabinData); err != nil {
		http.Error(w, "Error al actualizar la cabaña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cabaña actualizada con éxito"})
}
*/

func (h *AdminHandler) UpdateCabinHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	// 2. Obtener ID (generalmente viene en un input hidden o en el query)
	// AQUI HABRIA QUE MOFICIAR CON EL ID PROVENIENTE DEL LOGIN
	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "ID de cabaña requerido", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	cabinData := domain.Cabin{
		ID:           id,
		EmailContact: r.FormValue("email_contact"),
		PhoneContact: r.FormValue("phone_contact"),
		Password:     r.FormValue("password"),
	}

	if err := h.CabinServiceADM.UpdateCabin(&cabinData); err != nil {
		http.Error(w, "Error actualizando cabaña: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// HABRIA QUE REVEER LA REDIRECCION HACIA DONDE SE DEBERIA DIRIGIR
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
func (h *AdminHandler) GetAllReservationsHandler(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.ReservationServiceADM.GetAllReservations()
	if err != nil {
		http.Error(w, "Error al obtener las reservas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservations)
}
*/

func (h *AdminHandler) GetAllReservationsHandler(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.ReservationServiceADM.GetAllReservations()
	if err != nil {
		http.Error(w, "Error cargando reservas", http.StatusInternalServerError)
		return
	}

	component := views.Layout(reservations)
	component.Render(r.Context(), w)
}

/*
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
*/

/*
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

		if err := h.ReservationServiceUser.DeleteReservationByDate(cabinID, fecha); err != nil {
			http.Error(w, "Error al eliminar la reserva: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Reserva eliminada con éxito"})
	}
*/
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

	if err := h.ReservationServiceUser.DeleteReservationByDate(cabinID, fecha); err != nil {
		http.Error(w, "Error al eliminar la reserva: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	reservation, err := h.ReservationServiceUser.GetReservationByDate(fecha)
	if err != nil {
		http.Error(w, "Reserva no encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reservation)
}
