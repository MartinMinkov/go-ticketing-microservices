package events

type OrderCancelledEventTicketData struct {
	Id string `json:"id"`
}

type OrderCancelledEventData struct {
	Id      string `json:"id"`
	Version int64  `json:"version"`
	Ticket  OrderCancelledEventTicketData
}

type OrderCancelledEvent struct {
	Subject Subjects                `json:"subject"`
	Data    OrderCancelledEventData `json:"data"`
}

func NewOrderCancelledEvent(id string, ticketId string, version int64) *OrderCancelledEvent {
	return &OrderCancelledEvent{
		Subject: OrderCancelled,
		Data: OrderCancelledEventData{
			Id:      id,
			Version: version,
			Ticket: OrderCancelledEventTicketData{
				Id: ticketId,
			},
		},
	}
}
