package events

type TicketCreatedEventData struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Title  string `json:"title"`
	Price  int64  `json:"price"`
}

type TicketCreatedEvent struct {
	Subject Subjects               `json:"subject"`
	Data    TicketCreatedEventData `json:"data"`
}

func NewTicketCreatedEvent(id string, userId string, title string, price int64) *TicketCreatedEvent {
	return &TicketCreatedEvent{
		Subject: TicketCreated,
		Data: TicketCreatedEventData{
			Id:     id,
			UserId: userId,
			Title:  title,
			Price:  price,
		},
	}
}
