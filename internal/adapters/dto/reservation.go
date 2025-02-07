package dto

type ReservationRequest struct {
	NumCustomers int `json:"num_customers"`
}

type ReservationResponse struct {
	BookingId       string `json:"booking_id"`
	TablesReserved  int    `json:"tables_reserved"`
	RemainingTables int    `json:"remaining_tables"`
}

type CancelReservationResponse struct {
	FreedTables     int `json:"freed_tables"`
	RemainingTables int `json:"remaining_tables"`
}
