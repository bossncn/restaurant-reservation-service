package service_test

import (
	"errors"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	mockRepository "github.com/bossncn/restaurant-reservation-service/internal/core/repository/mock"
	"github.com/bossncn/restaurant-reservation-service/internal/core/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func TestReservationService(t *testing.T) {
	t.Run("NewReservationService", func(t *testing.T) {
		eventRequest := make(chan model.EventRequest, 100)
		logger := zap.NewNop()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mockRepository.NewMockReservationRepository(ctrl)

		svc := service.NewReservationService(mockRepo, logger, &eventRequest)

		assert.NotNil(t, svc)
	})
	t.Run("ReserveTables", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockReservationRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			svc := service.NewReservationService(mockRepo, logger, &eventRequest)

			// Mock event processor
			go func() {
				for req := range eventRequest {
					if req.Action == "reserve" {
						req.Response <- uuid.New().String()
					}
				}
			}()

			// Test reservation
			resID, numTables, err := svc.ReserveTables(6)

			assert.NoError(t, err)
			assert.NotEmpty(t, resID)
			assert.Equal(t, 2, numTables) // 6 customers require 2 tables
		})
		t.Run("InvalidCustomers", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockReservationRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			svc := service.NewReservationService(mockRepo, logger, &eventRequest)

			// Test with invalid customer count
			resID, numTables, err := svc.ReserveTables(0)

			assert.Error(t, err)
			assert.Equal(t, "number of customers must be greater than zero", err.Error())
			assert.Empty(t, resID)
			assert.Equal(t, 0, numTables)
		})
		t.Run("ErrorFromProcessor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockReservationRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			svc := service.NewReservationService(mockRepo, logger, &eventRequest)

			// Mock event processor
			go func() {
				for req := range eventRequest {
					if req.Action == "reserve" {
						req.Response <- errors.New("reservation failed")
					}
				}
			}()

			// Test reservation
			resID, numTables, err := svc.ReserveTables(6)

			assert.Error(t, err)
			assert.Equal(t, "reservation failed", err.Error())
			assert.Empty(t, resID)
			assert.Equal(t, 0, numTables)
		})
	})
	t.Run("CancelReservation", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockReservationRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			svc := service.NewReservationService(mockRepo, logger, &eventRequest)

			// Mock event processor
			go func() {
				for req := range eventRequest {
					if req.Action == "cancel" {
						req.Response <- 2 // Mock returning 2 tables freed
					}
				}
			}()

			// Test cancellation
			numTables, err := svc.CancelReservation("res-1")

			assert.NoError(t, err)
			assert.Equal(t, 2, numTables)
		})
		t.Run("ErrorFromProcessor", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mockRepository.NewMockReservationRepository(ctrl)
			eventRequest := make(chan model.EventRequest, 100)
			logger := zap.NewNop()
			svc := service.NewReservationService(mockRepo, logger, &eventRequest)

			// Mock event processor
			go func() {
				for req := range eventRequest {
					if req.Action == "cancel" {
						req.Response <- errors.New("cancellation failed")
					}
				}
			}()

			// Test cancellation
			numTables, err := svc.CancelReservation("res-1")

			assert.Error(t, err)
			assert.Equal(t, "cancellation failed", err.Error())
			assert.Equal(t, 0, numTables)
		})
	})
}
