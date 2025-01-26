package http

import (
	"errors"
	"github.com/bossncn/go-common/http/echo/response"
	"github.com/bossncn/go-common/http/model/error_code"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/dto"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TableHandler struct {
	Logger *zap.Logger
}

func NewTableHandler(logger *zap.Logger) *TableHandler {
	return &TableHandler{
		Logger: logger,
	}
}

// InitializeTable
// @Summary Initialize tables in the restaurant
// @Description Initializes the total number of tables in the restaurant. This endpoint must be called first and only once.
// @Tags table
// @Accept json
// @Produce json
// @Param request body dto.InitializeTableRequest true "Initialize Number of Table Request"
// @Success 200 {object} model.Response{data=dto.InitializeTableResponse} "Total Initialized Tables"
// @Failure 400 {object} model.Response{} "Table Already Initialized"
// @Router /public/table/init [post]
func (handler *TableHandler) InitializeTable(ctx echo.Context) error {
	var req dto.InitializeTableRequest
	if err := ctx.Bind(&req); err != nil {
		handler.Logger.Error("Failed to bind request", zap.Error(err))
		return response.Response(ctx, nil, errors.New(error_code.InvalidRequest))
	}
	return response.Response(ctx, dto.InitializeTableResponse{TotalTables: 0}, nil)
}
