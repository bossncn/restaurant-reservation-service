package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bossncn/go-common/http/model"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/dto"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func initializeTables(t *testing.T, echoInstance *echo.Echo, numTables int) {
	// Setup
	expected, _ := json.Marshal(model.CreateResponse("0", "Success", dto.InitializeTableResponse{TotalTables: numTables}))
	reqBody := fmt.Sprintf(`{"num_tables": %d}`, numTables)
	req := httptest.NewRequest(http.MethodPost, "/public/table/init", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Action
	echoInstance.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, string(expected), strings.TrimSpace(rec.Body.String()))
}

func TestIntegrationTables(t *testing.T) {
	t.Run("should return 200", func(t *testing.T) {
		echoInstance := Setup()

		initializeTables(t, echoInstance, 10)
	})
}
