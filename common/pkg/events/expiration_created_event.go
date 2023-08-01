package events

type ExpirationCreatedEventData struct {
	OrderId string `json:"order_id"`
}

type ExpirationCreatedEvent struct {
	Subject Subjects                   `json:"subject"`
	Data    ExpirationCreatedEventData `json:"data"`
}

func NewExpirationCreatedEvent(orderId string) *ExpirationCreatedEvent {
	return &ExpirationCreatedEvent{
		Subject: ExpirationCreated,
		Data: ExpirationCreatedEventData{
			OrderId: orderId,
		},
	}
}
