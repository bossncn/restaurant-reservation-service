package dto

type TablesResponse struct {
	TotalTables     int            `json:"total_tables"`
	AvailableTables int            `json:"available_tables"`
	Reservations    map[string]int `json:"reservations"`
}

type InitializeTableRequest struct {
	NumTables int `json:"num_tables"`
}

type InitializeTableResponse struct {
	TotalTables int `json:"total_tables"`
}
