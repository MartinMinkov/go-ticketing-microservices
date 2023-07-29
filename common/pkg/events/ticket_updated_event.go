package events

type TicketUpdatedEventData struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	OrderId string `json:"order_id"`
	Title   string `json:"title"`
	Price   int64  `json:"price"`
	Version int64  `json:"version"`
}

type TicketUpdatedEvent struct {
	Subject Subjects               `json:"subject"`
	Data    TicketUpdatedEventData `json:"data"`
}

func NewTicketUpdatedEvent(id string, userId string, orderId string, title string, price int64, version int64) *TicketUpdatedEvent {
	return &TicketUpdatedEvent{
		Subject: TicketUpdated,
		Data: TicketUpdatedEventData{
			Id:      id,
			UserId:  userId,
			OrderId: orderId,
			Title:   title,
			Price:   price,
			Version: version,
		},
	}
}
