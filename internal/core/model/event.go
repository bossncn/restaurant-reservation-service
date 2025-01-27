package model

type EventRequest struct {
	Id        string
	Action    string
	NumTables int
	ResID     string
	Response  chan interface{}
}
