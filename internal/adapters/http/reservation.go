package http

import (
	"errors"
	"github.com/bossncn/go-common/http/echo/response"
	"github.com/bossncn/go-common/http/model"
	"github.com/bossncn/go-common/http/model/error_code"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/dto"
	"github.com/bossncn/restaurant-reservation-service/internal/core/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ReservationHandler struct {
	logger             *zap.Logger
	reservationService service.ReservationService
	tableService       service.TableService
}

func NewReservationHandler(logger *zap.Logger, service *Service) *ReservationHandler {
	return &ReservationHandler{
		logger:             logger,
		reservationService: service.ReservationService,
		tableService:       service.TableService,
	}
}

func (handler *ReservationHandler) RegisterRoutes(_ *echo.Group, secureRoute *echo.Group) {
	secureReservationGroup := secureRoute.Group("/reservations")
	secureReservationGroup.POST("", handler.Reserve)
	secureReservationGroup.DELETE("/:id", handler.CancelReservation)
}

// Reserve
// @Summary Reserve tables
// @Description Reserves tables for a group of customers.
// @Tags Reservation
// @Accept json
// @Produce json
// @Param request body dto.ReservationRequest true "Number of customers in the group."
// @Success 200 {object} model.Response{data=dto.ReservationResponse} "Tables reserved successfully."
// @Failure 400 {object} model.Response{} "Reservation error."
// @Router /secure/reservations [post]
func (handler *ReservationHandler) Reserve(ctx echo.Context) error {
	var req dto.ReservationRequest
	if err := ctx.Bind(&req); err != nil {
		handler.logger.Error("Failed to bind request", zap.Error(err))
		return response.Response(ctx, nil, errors.New(error_code.InvalidRequest))
	}

	reservationID, numTables, err := handler.reservationService.ReserveTables(req.NumCustomers)

	if err != nil {
		handler.logger.Error("Failed to reserve tables", zap.Error(err))
		return response.Response(ctx, model.CreateError(error_code.InvalidRequest, err.Error()), err)
	}

	return response.Response(
		ctx,
		dto.ReservationResponse{
			BookingId:       reservationID,
			TablesReserved:  numTables,
			RemainingTables: handler.tableService.AvailableTables()},
		nil)
}

// CancelReservation
// @Summary Cancel a reservation
// @Description Cancels a reservation and releases the reserved tables.
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path string true "The reservation ID to cancel."
// @Success 200 {object} model.Response{data=dto.CancelReservationResponse} "Reservation canceled successfully."
// @Failure 400 {object} model.Response{} "Cancellation error."
// @Router /secure/reservations/{id} [delete]
func (handler *ReservationHandler) CancelReservation(ctx echo.Context) error {
	reservationID := ctx.Param("id")

	if reservationID == "" {
		handler.logger.Error("Missing ReservationID")
		return response.Response(ctx, nil, errors.New(error_code.InvalidRequest))
	}

	freedTables, err := handler.reservationService.CancelReservation(reservationID)

	if err != nil {
		handler.logger.Error("Failed to cancel reservation", zap.Error(err))
		return response.Response(ctx, model.CreateError(error_code.InvalidRequest, err.Error()), err)
	}

	return response.Response(
		ctx,
		dto.CancelReservationResponse{
			FreedTables:     freedTables,
			RemainingTables: handler.tableService.AvailableTables()},
		nil)
}
