package memory

import (
	"errors"
	"fmt"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"time"
)

type ReservationRepository struct {
	Reservations map[string]model.Reservation
}

func NewReservationRepository() *ReservationRepository {
	repo := &ReservationRepository{
		Reservations: make(map[string]model.Reservation),
	}
	return repo
}

func (r *ReservationRepository) CreateReservation(numTables int) *model.Reservation {
	return &model.Reservation{
		Id:        generateID(),
		NumTables: numTables,
	}
}

func (r *ReservationRepository) FindReservationById(id string) (*model.Reservation, error) {
	res, existed := r.Reservations[id]
	if !existed {
		return nil, errors.New("reservation not found")
	}

	return &res, nil
}

func (r *ReservationRepository) CancelReservation(reservationID string) error {
	res, err := r.FindReservationById(reservationID)
	if err != nil {
		return err
	}

	delete(r.Reservations, res.Id)

	return nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
