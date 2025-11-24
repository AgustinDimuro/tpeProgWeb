package application

import (
	"time"
	"tpeProgWeb/domain"
)

// DayView representa la información visual de un día en el calendario
type DayView struct {
	DayNumber    int
	IsActive     bool // True si el día pertenece al mes actual
	Reservations []*domain.Reservation
}

// BuildCalendarGrid genera la matriz de 5 o 6 semanas para el HTML
func BuildCalendarGrid(year int, month time.Month, reservations []*domain.Reservation) [][]DayView {
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	var weeks [][]DayView
	var currentWeek []DayView

	// Calcular padding inicial (días vacíos antes del 1ro del mes)
	// time.Weekday: Domingo=0 ... Sábado=6. Ajustamos para Lunes=0.
	startWeekday := int(firstOfMonth.Weekday()) - 1
	if startWeekday < 0 {
		startWeekday = 6 // Domingo
	}

	// Rellenar días vacíos al inicio
	for i := 0; i < startWeekday; i++ {
		currentWeek = append(currentWeek, DayView{IsActive: false})
	}

	// Rellenar los días reales del mes
	for day := 1; day <= lastOfMonth.Day(); day++ {
		currentDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		// Filtrar reservas para este día específico
		var dayRes []*domain.Reservation
		for _, r := range reservations {
			// Comparamos Año, Mes y Día
			y, m, d := r.Fecha.Date()
			if y == currentDate.Year() && m == currentDate.Month() && d == currentDate.Day() {
				dayRes = append(dayRes, r)
			}
		}

		currentWeek = append(currentWeek, DayView{
			DayNumber:    day,
			IsActive:     true,
			Reservations: dayRes,
		})

		// Si la semana se llena (7 días), la guardamos y empezamos nueva
		if len(currentWeek) == 7 {
			weeks = append(weeks, currentWeek)
			currentWeek = []DayView{}
		}
	}

	// Rellenar días vacíos al final si quedó la semana incompleta
	if len(currentWeek) > 0 {
		for len(currentWeek) < 7 {
			currentWeek = append(currentWeek, DayView{IsActive: false})
		}
		weeks = append(weeks, currentWeek)
	}

	return weeks
}
