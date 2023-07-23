package events

type TicketUpdatedEventData struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
}

type TicketUpdatedEvent struct {
	Subject Subjects               `json:"subject"`
	Data    TicketUpdatedEventData `json:"data"`
}
