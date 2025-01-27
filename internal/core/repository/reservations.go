package repository

import "github.com/bossncn/restaurant-reservation-service/internal/core/model"

type ReservationRepository interface {
	CreateReservation(numTables int) *model.Reservation
	FindReservationById(id string) (*model.Reservation, error)
	CancelReservation(reservationID string) error
}
