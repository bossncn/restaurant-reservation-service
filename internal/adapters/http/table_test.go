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

func TestTableHandler(t *testing.T) {
	t.Run("InitializeTable", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewTableHandler(logger, &http.Service{TableService: mockTableService})

			// Set up Echo mock context
			reqBody := dto.InitializeTableRequest{NumTables: 10}
			reqJSON, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(netHttp.MethodPost, "/table/init", bytes.NewReader(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Mock behavior
			mockTableService.EXPECT().InitializeTables(reqBody.NumTables).Return(nil).Times(1)

			// Execute handler
			err := handler.InitializeTable(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusOK, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, float64(reqBody.NumTables), res.Data.(map[string]interface{})["total_tables"])
		})
		t.Run("BindError", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock dependencies
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewTableHandler(logger, &http.Service{TableService: mockTableService})

			// Set up Echo mock context with invalid JSON
			req := httptest.NewRequest(netHttp.MethodPost, "/table/init", bytes.NewReader([]byte("invalid json")))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Execute handler
			err := handler.InitializeTable(ctx)

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
			mockTableService := serviceMock.NewMockTableService(ctrl)
			logger := zap.NewNop()
			handler := http.NewTableHandler(logger, &http.Service{TableService: mockTableService})

			// Set up Echo mock context
			reqBody := dto.InitializeTableRequest{NumTables: 10}
			reqJSON, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(netHttp.MethodPost, "/table/init", bytes.NewReader(reqJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(req, rec)

			// Mock behavior
			mockTableService.EXPECT().InitializeTables(reqBody.NumTables).Return(errors.New("table already initialized")).Times(1)

			// Execute handler
			err := handler.InitializeTable(ctx)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, netHttp.StatusBadRequest, rec.Code)

			var res model.Response
			err = json.Unmarshal(rec.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, error_code.InvalidRequest, res.Code)
			assert.Equal(t, "table already initialized", res.Data)
		})
	})
}
