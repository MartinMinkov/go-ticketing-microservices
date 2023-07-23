package events

type TicketUpdatedEventData struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	Title  string `json:"title"`
	Price  int64  `json:"price"`
}

type TicketUpdatedEvent struct {
	Subject Subjects               `json:"subject"`
	Data    TicketUpdatedEventData `json:"data"`
}

func NewTicketUpdatedEvent(id string, userId string, title string, price int64) *TicketUpdatedEvent {
	return &TicketUpdatedEvent{
		Subject: TicketUpdated,
		Data: TicketUpdatedEventData{
			Id:     id,
			UserId: userId,
			Title:  title,
			Price:  price,
		},
	}
}
