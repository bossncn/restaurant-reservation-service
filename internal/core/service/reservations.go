package service

import (
	"errors"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"github.com/bossncn/restaurant-reservation-service/internal/core/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
)

type ReservationService interface {
	ReserveTables(numCustomers int) (string, int, error)
	CancelReservation(reservationID string) (int, error)
}

type ReservationServiceImpl struct {
	reservationRepo repository.ReservationRepository
	logger          *zap.Logger
	requests        chan model.EventRequest
	wg              sync.WaitGroup
}

func NewReservationService(repo repository.ReservationRepository, logger *zap.Logger, eventRequest *chan model.EventRequest) *ReservationServiceImpl {
	repositoryService := &ReservationServiceImpl{
		reservationRepo: repo,
		logger:          logger,
		requests:        *eventRequest,
	}
	repositoryService.wg.Add(1)
	return repositoryService
}

func (s *ReservationServiceImpl) ReserveTables(numCustomers int) (string, int, error) {
	if numCustomers <= 0 {
		return "", 0, errors.New("number of customers must be greater than zero")
	}
	numTables := (numCustomers + 3) / 4 // Calculate required tables

	resp := make(chan interface{})
	s.requests <- model.EventRequest{Id: (uuid.New()).String(), Action: "reserve", NumTables: numTables, Response: resp}
	result := <-resp

	if err, ok := result.(error); ok {
		return "", 0, err
	}

	return result.(string), numTables, nil
}

func (s *ReservationServiceImpl) CancelReservation(reservationID string) (int, error) {
	resp := make(chan interface{})
	s.requests <- model.EventRequest{Id: (uuid.New()).String(), Action: "cancel", ResID: reservationID, Response: resp}
	result := <-resp
	if err, ok := result.(error); ok {
		return 0, err
	}
	return result.(int), nil
}
