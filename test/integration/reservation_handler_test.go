package integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/bossncn/go-common/http/model"
	"github.com/bossncn/restaurant-reservation-service/internal/adapters/dto"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationReservations(t *testing.T) {
	t.Run("should return 200 OK", func(t *testing.T) {
		echoInstance := Setup()

		// Setup
		initializeTables(t, echoInstance, 2)

		reqBody := `{"num_customers": 5}`
		req := httptest.NewRequest(http.MethodPost, "/secure/reservations", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Action
		echoInstance.ServeHTTP(rec, req)

		// Assert
		var resp model.Response
		_ = json.Unmarshal([]byte(rec.Body.String()), &resp)

		jsonData, _ := json.Marshal(resp.Data)

		var data dto.ReservationResponse
		_ = json.Unmarshal(jsonData, &data)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotNil(t, data.BookingId)
		assert.Equal(t, 2, data.TablesReserved)
		assert.Equal(t, 0, data.RemainingTables)
	})
	t.Run("should return 400 Bad Request", func(t *testing.T) {
		echoInstance := Setup()

		reqBody := `{"num_customers": 5}`
		req := httptest.NewRequest(http.MethodPost, "/secure/reservations", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// Action
		echoInstance.ServeHTTP(rec, req)

		// Assert
		var resp model.Response
		_ = json.Unmarshal([]byte(rec.Body.String()), &resp)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "tables has not been initialized", resp.Data)
	})
}
