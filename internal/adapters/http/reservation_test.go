package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bossncn/go-common/http/model"
	"github.com/bossncn/go-common/http/model/error_code"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/dto"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/http"
	serviceMock "github.com/bossncn/restaurant-reservation-service/internal/core/service/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	netHttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestReservationHandler(t *testing.T) {
	t.Run("Reserve", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockReservationService := serviceMock.NewMockReservationService(ctrl)
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewReservationHandler(logger, &http.Service{
				ReservationService: mockReservationService,
				TableService:       mockTableService,
			})

			// Set up Echo mock context
			reqBody := dto.ReservationRequest{NumCustomers: 8}
			reqJSON, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(netHttp.MethodPost, "/reservations", bytes.NewReader(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Mock behavior
			mockReservationService.EXPECT().ReserveTables(reqBody.NumCustomers).Return("res-1", 2, nil).Times(1)
			mockTableService.EXPECT().AvailableTables().Return(8).Times(1)

			// Execute handler
			err := handler.Reserve(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusOK, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, "res-1", res.Data.(map[string]interface{})["booking_id"])
			assert.Equal(t, float64(2), res.Data.(map[string]interface{})["tables_reserved"])
			assert.Equal(t, float64(8), res.Data.(map[string]interface{})["remaining_tables"])
		})
		t.Run("BindError", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockReservationService := serviceMock.NewMockReservationService(ctrl)
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewReservationHandler(logger, &http.Service{
				ReservationService: mockReservationService,
				TableService:       mockTableService,
			})

			// Set up Echo mock context with invalid JSON
			req := httptest.NewRequest(netHttp.MethodPost, "/reservations", bytes.NewReader([]byte("invalid json")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Execute handler
			err := handler.Reserve(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusBadRequest, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, error_code.InvalidRequest, res.Code)
		})
		t.Run("ServiceError", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockReservationService := serviceMock.NewMockReservationService(ctrl)
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewReservationHandler(logger, &http.Service{
				ReservationService: mockReservationService,
				TableService:       mockTableService,
			})

			// Set up Echo mock context
			reqBody := dto.ReservationRequest{NumCustomers: 8}
			reqJSON, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(netHttp.MethodPost, "/reservations", bytes.NewReader(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Mock behavior
			mockReservationService.EXPECT().ReserveTables(reqBody.NumCustomers).Return("", 0, errors.New("reservation failed")).Times(1)

			// Execute handler
			err := handler.Reserve(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusBadRequest, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, error_code.InvalidRequest, res.Code)
			assert.Equal(t, "reservation failed", res.Data)
		})
	})
	t.Run("CancelReservation", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockReservationService := serviceMock.NewMockReservationService(ctrl)
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewReservationHandler(logger, &http.Service{
				ReservationService: mockReservationService,
				TableService:       mockTableService,
			})

			// Set up Echo mock context
			req := httptest.NewRequest(netHttp.MethodDelete, "/reservations/res-1", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("res-1")

			// Mock behavior
			mockReservationService.EXPECT().CancelReservation("res-1").Return(3, nil).Times(1)
			mockTableService.EXPECT().AvailableTables().Return(10).Times(1)

			// Execute handler
			err := handler.CancelReservation(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusOK, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, float64(3), res.Data.(map[string]interface{})["freed_tables"])
			assert.Equal(t, float64(10), res.Data.(map[string]interface{})["remaining_tables"])
		})
		t.Run("MissingID", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockReservationService := serviceMock.NewMockReservationService(ctrl)
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewReservationHandler(logger, &http.Service{
				ReservationService: mockReservationService,
				TableService:       mockTableService,
			})

			// Set up Echo mock context without ID
			req := httptest.NewRequest(netHttp.MethodDelete, "/reservations/", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Execute handler
			err := handler.CancelReservation(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusBadRequest, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, error_code.InvalidRequest, res.Code)
		})
		t.Run("ServiceError", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockReservationService := serviceMock.NewMockReservationService(ctrl)
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewReservationHandler(logger, &http.Service{
				ReservationService: mockReservationService,
				TableService:       mockTableService,
			})

			// Set up Echo mock context
			req := httptest.NewRequest(netHttp.MethodDelete, "/reservations/res-1", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues("res-1")

			// Mock behavior
			mockReservationService.EXPECT().CancelReservation("res-1").Return(0, errors.New("cancellation failed")).Times(1)

			// Execute handler
			err := handler.CancelReservation(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusBadRequest, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, error_code.InvalidRequest, res.Code)
			assert.Equal(t, "cancellation failed", res.Data)
		})
	})
}
