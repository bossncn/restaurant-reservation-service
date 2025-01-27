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

type TableHandler struct {
	logger       *zap.Logger
	tableService *service.TableService
}

func NewTableHandler(logger *zap.Logger, service *Service) *TableHandler {
	return &TableHandler{
		logger:       logger,
		tableService: service.TableService,
	}
}

func (handler *TableHandler) RegisterRoutes(publicRouteGroup *echo.Group) {
	publicTableRouteGroup := publicRouteGroup.Group("/table")
	publicTableRouteGroup.POST("/init", handler.InitializeTable)
}

// InitializeTable
// @Summary Initialize tables in the restaurant
// @Description Initializes the total number of tables in the restaurant. This endpoint must be called first and only once.
// @Tags table
// @Accept json
// @Produce json
// @Param request body dto.InitializeTableRequest true "Initialize Number of table Request"
// @Success 200 {object} model.Response{data=dto.InitializeTableResponse} "Total Initialized Tables"
// @Failure 400 {object} model.Response{} "table Already Initialized"
// @Router /public/table/init [post]
func (handler *TableHandler) InitializeTable(ctx echo.Context) error {
	var req dto.InitializeTableRequest
	if err := ctx.Bind(&req); err != nil {
		handler.logger.Error("Failed to bind request", zap.Error(err))
		return response.Response(ctx, nil, errors.New(error_code.InvalidRequest))
	}

	if err := handler.tableService.InitializeTables(req.NumTables); err != nil {
		handler.logger.Error("Failed to initialize tables", zap.Error(err))
		return response.Response(ctx, model.CreateError(error_code.InvalidRequest, err.Error()), err)
	}

	return response.Response(ctx, dto.InitializeTableResponse{TotalTables: req.NumTables}, nil)
}
