package memory_test

import (
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryReservationRepository(t *testing.T) {
	t.Run("NewReservationRepository", func(t *testing.T) {
		repo := memory.NewReservationRepository()
		assert.NotNil(t, repo)
		assert.Empty(t, repo.Reservations)
	})
	t.Run("CreateReservation", func(t *testing.T) {
		repo := memory.NewReservationRepository()

		// Create a reservation
		reservation := repo.CreateReservation(3)

		assert.NotNil(t, reservation)
		assert.Equal(t, 3, reservation.NumTables)
		assert.NotEmpty(t, reservation.Id)
	})
	t.Run("FindReservationById", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := memory.NewReservationRepository()

			// Create and add a reservation to the repository
			reservation := repo.CreateReservation(2)
			repo.Reservations[reservation.Id] = *reservation

			// Find the reservation
			found, err := repo.FindReservationById(reservation.Id)

			assert.NoError(t, err)
			assert.Equal(t, reservation.Id, found.Id)
			assert.Equal(t, reservation.NumTables, found.NumTables)
		})
		t.Run("NotFound", func(t *testing.T) {
			repo := memory.NewReservationRepository()

			// Try to find a non-existent reservation
			_, err := repo.FindReservationById("non-existent-id")

			assert.Error(t, err)
			assert.Equal(t, "reservation not found", err.Error())
		})
	})
	t.Run("CancelReservation", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := memory.NewReservationRepository()

			// Create and add a reservation to the repository
			reservation := repo.CreateReservation(4)
			repo.Reservations[reservation.Id] = *reservation

			// Cancel the reservation
			err := repo.CancelReservation(reservation.Id)

			assert.NoError(t, err)
			assert.NotContains(t, repo.Reservations, reservation.Id)
		})
		t.Run("NotFound", func(t *testing.T) {
			repo := memory.NewReservationRepository()

			// Try to cancel a non-existent reservation
			err := repo.CancelReservation("non-existent-id")

			assert.Error(t, err)
			assert.Equal(t, "reservation not found", err.Error())
		})
	})
}
