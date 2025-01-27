package memory_test

import (
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/memory"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryTableRepository(t *testing.T) {
	t.Run("NewTableRepository", func(t *testing.T) {
		repo := memory.NewTableRepository()

		assert.NotNil(t, repo)
		assert.Equal(t, 0, repo.AvailableTables())
		assert.False(t, repo.IsTableInitialized())
	})
	t.Run("InitializeTables", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := memory.NewTableRepository()

			err := repo.InitializeTables(10)

			assert.NoError(t, err)
			assert.Equal(t, 10, repo.AvailableTables())
			assert.True(t, repo.IsTableInitialized())
		})
		t.Run("AlreadyInitialized", func(t *testing.T) {
			repo := memory.NewTableRepository()

			_ = repo.InitializeTables(10)
			err := repo.InitializeTables(5)

			assert.Error(t, err)
			assert.Equal(t, "tables already initialized", err.Error())
			assert.Equal(t, 10, repo.AvailableTables())
		})
	})
	t.Run("ReserveTables", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := memory.NewTableRepository()

			_ = repo.InitializeTables(10)
			reservation := model.Reservation{
				Id:        "res-1",
				NumTables: 3,
			}

			err := repo.ReserveTables(reservation)

			assert.NoError(t, err)
			assert.Equal(t, 7, repo.AvailableTables())
			assert.Contains(t, repo.Table.Reservations, "res-1")
		})
	})
	t.Run("CancelReservedTable", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			repo := memory.NewTableRepository()

			_ = repo.InitializeTables(10)
			reservation := model.Reservation{
				Id:        "res-1",
				NumTables: 3,
			}
			_ = repo.ReserveTables(reservation)

			err := repo.CancelReservedTable("res-1")

			assert.NoError(t, err)
			assert.Equal(t, 10, repo.AvailableTables())
			assert.NotContains(t, repo.Table.Reservations, "res-1")
		})
		t.Run("NotFound", func(t *testing.T) {
			repo := memory.NewTableRepository()

			_ = repo.InitializeTables(10)

			err := repo.CancelReservedTable("non-existent-id")

			assert.Error(t, err)
			assert.Equal(t, "booking not found", err.Error())
			assert.Equal(t, 10, repo.AvailableTables())
		})
	})
	t.Run("AvailableTables", func(t *testing.T) {
		repo := memory.NewTableRepository()

		_ = repo.InitializeTables(10)
		reservation := model.Reservation{
			Id:        "res-1",
			NumTables: 4,
		}
		_ = repo.ReserveTables(reservation)

		assert.Equal(t, 6, repo.AvailableTables())
	})
}
