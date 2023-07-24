package events

type OrderCancelledEventTicketData struct {
	Id string `json:"id"`
}

type OrderCancelledEventData struct {
	Id     string `json:"id"`
	ticket OrderCancelledEventTicketData
}

type OrderCancelledEvent struct {
	Subject Subjects                `json:"subject"`
	Data    OrderCancelledEventData `json:"data"`
}

func NewOrderCancelledEvent(id string, ticketId string) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		Subject: OrderCancelled,
		Data: OrderCancelledEventData{
			Id: id,
			ticket: OrderCancelledEventTicketData{
				Id: ticketId,
			},
		},
	}
}
