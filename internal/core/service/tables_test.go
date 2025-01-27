package service_test

import (
	"errors"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	mockRepository "github.com/bossncn/restaurant-reservation-service/internal/core/repository/mock"
	"github.com/bossncn/restaurant-reservation-service/internal/core/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func TestTablesService(t *testing.T) {
	t.Run("NewTableService", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockRepository.NewMockTableRepository(ctrl)
		eventRequest := make(chan model.EventRequest, 100)
		logger := zap.NewNop()

		tableService := service.NewTableService(mockRepo, logger, &eventRequest)

		assert.NotNil(t, tableService)
	})
	t.Run("InitializeTables", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockTableRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			tableService := service.NewTableService(mockRepo, logger, &eventRequest)

			// Mock event processor
			go func() {
				for req := range eventRequest {
					if req.Action == "initialize" {
						req.Response <- nil // Simulate success
					}
				}
			}()

			// Test initialization
			err := tableService.InitializeTables(10)

			assert.NoError(t, err)
		})
		t.Run("Error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockTableRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			tableService := service.NewTableService(mockRepo, logger, &eventRequest)

			// Mock event processor
			go func() {
				for req := range eventRequest {
					if req.Action == "initialize" {
						req.Response <- errors.New("initialization failed") // Simulate failure
					}
				}
			}()

			// Test initialization
			err := tableService.InitializeTables(10)

			assert.Error(t, err)
			assert.Equal(t, "initialization failed", err.Error())
		})
	})
	t.Run("AvailableTables", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockRepository.NewMockTableRepository(ctrl)
		eventRequest := make(chan model.EventRequest, 100)
		logger := zap.NewNop()
		tableService := service.NewTableService(mockRepo, logger, &eventRequest)

		// Mock AvailableTables behavior
		mockRepo.EXPECT().AvailableTables().Return(8).Times(1)

		// Test available tables
		availableTables := tableService.AvailableTables()

		assert.Equal(t, 8, availableTables)
	})
}
