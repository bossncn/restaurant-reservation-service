package http

import (
	"errors"
	"github.com/bossncn/go-common/http/echo/response"
	"github.com/bossncn/go-common/http/model/error_code"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/dto"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ReservationHandler struct {
	Logger *zap.Logger
}

func NewReservationHandler(logger *zap.Logger) *ReservationHandler {
	return &ReservationHandler{
		Logger: logger,
	}
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
// @Router /public/reservations [post]
func (handler *ReservationHandler) Reserve(ctx echo.Context) error {
	var req dto.ReservationRequest
	if err := ctx.Bind(&req); err != nil {
		handler.Logger.Error("Failed to bind request", zap.Error(err))
		return response.Response(ctx, nil, errors.New(error_code.InvalidRequest))
	}
	return response.Response(
		ctx,
		dto.ReservationResponse{
			Id:              "",
			TablesReserved:  0,
			RemainingTables: 0},
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
// @Router /public/reservations/{id} [delete]
func (handler *ReservationHandler) CancelReservation(ctx echo.Context) error {
	reservationID := ctx.Param("id")

	if reservationID == "" {
		handler.Logger.Error("Missing ReservationID")
		return response.Response(ctx, nil, errors.New(error_code.InvalidRequest))
	}

	return response.Response(
		ctx,
		dto.CancelReservationResponse{
			FreedTables:     0,
			RemainingTables: 0},
		nil)
}
