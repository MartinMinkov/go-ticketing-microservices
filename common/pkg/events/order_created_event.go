package events

import "time"

type OrderCreatedEventTicketData struct {
	Id    string `json:"id"`
	Price int64  `json:"price"`
}

type OrderCreatedEventData struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Status    string    `json:"title"`
	ExpiresAt time.Time `json:"expires_at"`
	ticket    OrderCreatedEventTicketData
}

type OrderCreatedEvent struct {
	Subject Subjects              `json:"subject"`
	Data    OrderCreatedEventData `json:"data"`
}

func NewOrderCreatedEvent(id string, userId string, ticketId string, status string, price int64, expiresAt time.Time) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		Subject: OrderCreated,
		Data: OrderCreatedEventData{
			Id:        id,
			UserId:    userId,
			Status:    status,
			ExpiresAt: expiresAt,
			ticket: OrderCreatedEventTicketData{
				Id:    ticketId,
				Price: price,
			},
		},
	}
}
