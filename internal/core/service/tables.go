package service

import (
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"github.com/bossncn/restaurant-reservation-service/internal/core/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
)

type TableService interface {
	InitializeTables(numTables int) error
	AvailableTables() int
}

type TableServiceImpl struct {
	tableRepo repository.TableRepository
	requests  chan model.EventRequest
	logger    *zap.Logger
	wg        sync.WaitGroup
}

func NewTableService(repo repository.TableRepository, logger *zap.Logger, eventRequest *chan model.EventRequest) *TableServiceImpl {
	tableService := &TableServiceImpl{
		tableRepo: repo,
		logger:    logger,
		requests:  *eventRequest,
	}
	tableService.wg.Add(1)
	return tableService
}

func (s *TableServiceImpl) InitializeTables(numTables int) error {
	resp := make(chan interface{})
	s.requests <- model.EventRequest{Id: (uuid.New()).String(), Action: "initialize", NumTables: numTables, Response: resp}
	result := <-resp
	if err, ok := result.(error); ok {
		return err
	}
	return nil
}

func (s *TableServiceImpl) AvailableTables() int {
	return s.tableRepo.AvailableTables()
}
