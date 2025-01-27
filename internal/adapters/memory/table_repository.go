package memory

import (
	"errors"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
)

type TableRepository struct {
	Table model.Table
}

func NewTableRepository() *TableRepository {
	repo := &TableRepository{
		Table: model.Table{
			Reservations: make(map[string]model.Reservation),
		},
	}
	return repo
}

func (r *TableRepository) InitializeTables(numTables int) error {
	if r.Table.TotalTables > 0 {
		return errors.New("tables already initialized")
	}

	r.Table.TotalTables = numTables
	r.Table.AvailableTables = numTables
	return nil
}

func (r *TableRepository) ReserveTables(reservation model.Reservation) error {
	r.Table.Reservations[reservation.Id] = reservation
	r.Table.AvailableTables -= reservation.NumTables

	return nil
}

func (r *TableRepository) CancelReservedTable(reservationId string) error {
	reservation, existed := r.Table.Reservations[reservationId]

	if !existed {
		return errors.New("booking not found")
	}
	delete(r.Table.Reservations, reservation.Id)
	r.Table.AvailableTables += reservation.NumTables

	return nil
}

func (r *TableRepository) IsTableInitialized() bool {
	return r.Table.TotalTables > 0
}

func (r *TableRepository) AvailableTables() int {
	return r.Table.AvailableTables
}
