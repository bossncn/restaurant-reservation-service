package repository

import "github.com/bossncn/restaurant-reservation-service/internal/core/model"

type TableRepository interface {
	InitializeTables(numTables int) error
	ReserveTables(reservation model.Reservation) error
	CancelReservedTable(reservationId string) error
	AvailableTables() int
	IsTableInitialized() bool
}
