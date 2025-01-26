package dto

import "github.com/bossncn/restaurant-reservation-service/internal/core/model"

type TablesResponse struct {
	TotalTables     int                    `json:"total_tables"`
	AvailableTables int                    `json:"available_tables"`
	Reservations    map[string]model.Table `json:"reservations"`
}

type InitializeTableRequest struct {
	NumTables int `json:"num_tables"`
}

type InitializeTableResponse struct {
	TotalTables int `json:"total_tables"`
}
