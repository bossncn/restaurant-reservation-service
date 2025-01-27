package event

import (
	"errors"
	"fmt"
	"github.com/bossncn/restaurant-reservation-service/internal/core/model"
	"github.com/bossncn/restaurant-reservation-service/internal/core/repository"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Processor struct {
	tableRepo       repository.TableRepository
	reservationRepo repository.ReservationRepository
	requests        chan model.EventRequest
	stopChan        chan bool
	wg              sync.WaitGroup
	logger          *zap.Logger
}

func NewProcessor(tableRepository repository.TableRepository, reservationRepository repository.ReservationRepository, logger *zap.Logger) (*Processor, *chan model.EventRequest) {
	requests := make(chan model.EventRequest, 100)

	processor := &Processor{
		tableRepo:       tableRepository,
		reservationRepo: reservationRepository,
		requests:        requests,
		stopChan:        make(chan bool),
		logger:          logger,
	}

	processor.wg.Add(1)
	return processor, &requests
}

func (e *Processor) ProcessRequests() {
	defer e.wg.Done()
	for {
		select {
		case req := <-e.requests:
			timeStarted := time.Now()
			e.logger.Info("Incoming Event EventRequest", zap.String("requestId", req.Id), zap.String("action", req.Action))
			switch req.Action {
			case "initialize":
				err := e.tableRepo.InitializeTables(req.NumTables)
				if err != nil {
					req.Response <- e.logError(req.Id, "initialize", err)
				} else {
					req.Response <- err
				}
			case "reserve":
				if !e.tableRepo.IsTableInitialized() {
					req.Response <- e.logError(req.Id, "reserve", errors.New("tables has not been initialized"))
				} else if req.NumTables <= 0 {
					req.Response <- e.logError(req.Id, "reserve", errors.New("invalid number of tables"))
				} else if req.NumTables > e.tableRepo.AvailableTables() {
					req.Response <- e.logError(req.Id, "reserve", errors.New("not enough tables available"))
				} else {
					reservation := e.reservationRepo.CreateReservation(req.NumTables)
					err := e.tableRepo.ReserveTables(*reservation)
					if err != nil {
						req.Response <- err
					}
					req.Response <- reservation.Id
				}
			case "cancel":
				if !e.tableRepo.IsTableInitialized() {
					req.Response <- e.logError(req.Id, "cancel", errors.New("tables has not been initialized"))
				} else {
					reservation, err := e.reservationRepo.FindReservationById(req.ResID)
					if err != nil {
						req.Response <- e.logError(req.Id, "cancel", err)
					} else {
						err := e.tableRepo.CancelReservedTable(reservation.Id)
						if err != nil {
							req.Response <- e.logError(req.Id, "cancel", err)
						} else {
							err = e.reservationRepo.CancelReservation(reservation.Id)
							if err != nil {
								req.Response <- e.logError(req.Id, "cancel", err)
							}
							req.Response <- reservation.NumTables
						}
					}
				}
			}
			e.logger.Info("Event EventRequest Complete", zap.String("requestId", req.Id), zap.String("action", req.Action), zap.String("elapsed", fmt.Sprintf("%.3f ms", float64(time.Since(timeStarted).Microseconds())/1000)))
		case <-e.stopChan:
			return
		}
	}
}

func (e *Processor) logError(requestId string, action string, err error) error {
	switch action {
	case "initialize":
		e.logger.Error("Error Initialize tables", zap.String("requestId", requestId), zap.Error(err))
		return err
	case "reserve":
		e.logger.Error("Error Reserve tables", zap.String("requestId", requestId), zap.Error(err))
		return err
	case "cancel":
		e.logger.Error("Error Cancel tables", zap.String("requestId", requestId), zap.Error(err))
		return err
	default:
		e.logger.Error("Error Process Event", zap.String("requestId", requestId), zap.Error(err))
		return err
	}
}
