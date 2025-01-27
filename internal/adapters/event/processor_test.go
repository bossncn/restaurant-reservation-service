package event_test

import (
	"errors"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/event"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	mockRepository "github.com/bossncn/restaurant-reservation-service/internal/core/repository/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEventProcessor_Initialize(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().InitializeTables(10).Return(nil).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-1",
			Action:    "initialize",
			NumTables: 10,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.Nil(t, res)
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().InitializeTables(10).Return(errors.New("already initialized")).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-2",
			Action:    "initialize",
			NumTables: 10,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "already initialized")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
}

func TestEventProcessor_Reserve(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockTableRepo.EXPECT().AvailableTables().Return(5).Times(1)
		mockReservationRepo.EXPECT().CreateReservation(3).Return(&model.Reservation{Id: "res-1", NumTables: 3}).Times(1)
		mockTableRepo.EXPECT().ReserveTables(model.Reservation{Id: "res-1", NumTables: 3}).Return(nil).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-3",
			Action:    "reserve",
			NumTables: 3,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.Equal(t, "res-1", res)
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("TableNotInitialized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(false).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-4",
			Action:    "reserve",
			NumTables: 3,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "tables has not been initialized")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("InvalidNumberOfTable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-4",
			Action:    "reserve",
			NumTables: 0,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "invalid number of tables")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("NotEnoughTable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockTableRepo.EXPECT().AvailableTables().Return(5).Times(1)
		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-3",
			Action:    "reserve",
			NumTables: 6,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "not enough tables available")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("ErrorReserveTable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockTableRepo.EXPECT().AvailableTables().Return(5).Times(1)
		mockReservationRepo.EXPECT().CreateReservation(3).Return(&model.Reservation{Id: "res-1", NumTables: 3}).Times(1)
		mockTableRepo.EXPECT().ReserveTables(model.Reservation{Id: "res-1", NumTables: 3}).Return(errors.New("something went wrong")).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:        "req-3",
			Action:    "reserve",
			NumTables: 3,
			Response:  response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "something went wrong")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
}

func TestEventProcessor_Cancel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockReservationRepo.EXPECT().FindReservationById("res-1").Return(&model.Reservation{Id: "res-1", NumTables: 3}, nil).Times(1)
		mockTableRepo.EXPECT().CancelReservedTable("res-1").Return(nil).Times(1)
		mockReservationRepo.EXPECT().CancelReservation("res-1").Return(nil).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:       "req-5",
			Action:   "cancel",
			ResID:    "res-1",
			Response: response,
		}

		select {
		case res := <-response:
			assert.Equal(t, 3, res)
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("TableNotInitialized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(false).Times(1)
		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:       "req-5",
			Action:   "cancel",
			ResID:    "res-1",
			Response: response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "tables has not been initialized")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("ReservationIdNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockReservationRepo.EXPECT().FindReservationById("res-1").Return(nil, errors.New("not found")).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:       "req-5",
			Action:   "cancel",
			ResID:    "res-1",
			Response: response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "not found")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("CancelReservedTableFailed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockReservationRepo.EXPECT().FindReservationById("res-1").Return(&model.Reservation{Id: "res-1", NumTables: 3}, nil).Times(1)
		mockTableRepo.EXPECT().CancelReservedTable("res-1").Return(errors.New("something went wrong")).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:       "req-5",
			Action:   "cancel",
			ResID:    "res-1",
			Response: response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "something went wrong")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
	t.Run("CancelReservationFailed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
		mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
		logger := zap.NewNop()

		processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

		mockTableRepo.EXPECT().IsTableInitialized().Return(true).Times(1)
		mockReservationRepo.EXPECT().FindReservationById("res-1").Return(&model.Reservation{Id: "res-1", NumTables: 3}, nil).Times(1)
		mockTableRepo.EXPECT().CancelReservedTable("res-1").Return(nil).Times(1)
		mockReservationRepo.EXPECT().CancelReservation("res-1").Return(errors.New("something went wrong")).Times(1)

		go processor.ProcessRequests()

		response := make(chan interface{}, 1)
		*requests <- model.EventRequest{
			Id:       "req-5",
			Action:   "cancel",
			ResID:    "res-1",
			Response: response,
		}

		select {
		case res := <-response:
			assert.EqualError(t, res.(error), "something went wrong")
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for response")
		}
	})
}

func TestEventProcessor_InvalidAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTableRepo := mockRepository.NewMockTableRepository(ctrl)
	mockReservationRepo := mockRepository.NewMockReservationRepository(ctrl)
	logger := zap.NewNop()

	processor, requests := event.NewProcessor(mockTableRepo, mockReservationRepo, logger)

	go processor.ProcessRequests()

	response := make(chan interface{}, 1)
	*requests <- model.EventRequest{
		Id:       "req-6",
		Action:   "unknown",
		Response: response,
	}

	select {
	case <-response:
		t.Fatal("unexpected response for invalid action")
	case <-time.After(1 * time.Second):
		// Expected timeout for invalid action
	}
}
