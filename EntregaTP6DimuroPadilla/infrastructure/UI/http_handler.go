package ui

import (
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
	// Obtenemos todas las reservaciones y las pasamos al Layout
	reservations, err := h.ReservationServiceUser.GetAllReservations()
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
	fechaStr := r.FormValue("fecha") // El input type="date" envía "2025-01-25"

	cabinID, err := strconv.ParseInt(cabinIDStr, 10, 64)
	if err != nil {
		http.Error(w, "El 'cabin_id' debe ser un número", http.StatusBadRequest)
		return
	}

	// Parseamos la fecha
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

	reservations, err := h.ReservationServiceUser.GetAllReservations()
	if err != nil {
		http.Error(w, "Error al recargar la lista", http.StatusInternalServerError)
		return
	}

	views.ReservationList(reservations).Render(r.Context(), w)
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

	// 1. Parsear fecha vieja
	oldDate, err := time.Parse("2006-01-02", oldDateStr)
	if err != nil {
		http.Error(w, "El formato de 'fecha' (actual) debe ser AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	// 2. Parsear fecha nueva
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

	reservations, err := h.ReservationServiceUser.GetAllReservations()
	if err != nil {
		http.Error(w, "Error al recargar la lista", http.StatusInternalServerError)
		return
	}

	views.ReservationList(reservations).Render(r.Context(), w)
}

func (h *UserHandler) HandleShowCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Leer parámetros de la URL
	queryMonth := r.URL.Query().Get("month")
	queryYear := r.URL.Query().Get("year")

	if queryMonth != "" && queryYear != "" {
		m, errM := strconv.Atoi(queryMonth)
		y, errY := strconv.Atoi(queryYear)
		if errM == nil && errY == nil && m >= 1 && m <= 12 {
			currentMonth = m
			currentYear = y
		}
	}

	// Calcular lógica de navegación (Anterior / Siguiente)
	// Creamos una fecha base con el mes actual visualizado
	targetDate := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, time.UTC)

	prevDate := targetDate.AddDate(0, -1, 0) // Restamos 1 mes
	nextDate := targetDate.AddDate(0, 1, 0)  // Sumamos 1 mes

	// Obtener reservas (Igual que antes)
	reservations, err := h.ReservationServiceUser.GetAllReservations()
	if err != nil {
		http.Error(w, "Error al cargar reservas", http.StatusInternalServerError)
		return
	}

	// Construir la grilla
	weeks := application.BuildCalendarGrid(currentYear, time.Month(currentMonth), reservations)

	// Renderizar
	// Pasamos más datos: mes/año actual, y mes/año de navegación
	nombresMeses := []string{"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

	component := views.CalendarPage(
		nombresMeses[currentMonth], // Nombre visual
		currentMonth,               // Número del mes actual (para el selector)
		currentYear,                // Año actual
		weeks,
		int(prevDate.Month()), prevDate.Year(), // Datos para flecha izquierda
		int(nextDate.Month()), nextDate.Year(), // Datos para flecha derecha
	)
	component.Render(r.Context(), w)
}

func (h *AdminHandler) UpdateCabinHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	// Obtener ID (generalmente viene en un input hidden o en el query)
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

func (h *AdminHandler) GetAllReservationsHandler(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.ReservationServiceADM.GetAllReservations()
	if err != nil {
		http.Error(w, "Error cargando reservas", http.StatusInternalServerError)
		return
	}

	component := views.Layout(reservations)
	component.Render(r.Context(), w)
}

func (h *UserHandler) DeleteReservationHandler(w http.ResponseWriter, r *http.Request) {
	// Parsear el formulario
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	// Obtener valores
	cabinIDStr := r.FormValue("cabin_id")
	dateStr := r.FormValue("fecha")

	if cabinIDStr == "" || dateStr == "" {
		http.Error(w, "Se requieren los parámetros 'cabin_id' y 'fecha'", http.StatusBadRequest)
		return
	}

	// Conversiones
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

	// Llamar al servicio
	if err := h.ReservationServiceUser.DeleteReservationByDate(cabinID, fecha); err != nil {
		http.Error(w, "Error al eliminar la reserva: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	reservations := []*domain.Reservation{reservation}
	views.ReservationList(reservations).Render(r.Context(), w)
}
