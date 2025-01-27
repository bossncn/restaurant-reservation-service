package model

type Table struct {
	TotalTables     int
	AvailableTables int
	Reservations    map[string]Reservation
}
