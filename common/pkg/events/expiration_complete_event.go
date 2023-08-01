package events

type ExpirationCompleteEventData struct {
	OrderId string `json:"order_id"`
}

type ExpirationCompleteEvent struct {
	Subject Subjects                    `json:"subject"`
	Data    ExpirationCompleteEventData `json:"data"`
}

func NewExpirationCompleteEvent(orderId string) *ExpirationCompleteEvent {
	return &ExpirationCompleteEvent{
		Subject: ExpirationComplete,
		Data: ExpirationCompleteEventData{
			OrderId: orderId,
		},
	}
}
