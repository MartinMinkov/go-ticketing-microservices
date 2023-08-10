package events

type PaymentCreatedEventData struct {
	Id       string `json:"id"`
	OrderId  string `json:"order_id"`
	StripeId string `json:"stripe_id"`
}

type PaymentCreatedEvent struct {
	Subject Subjects                `json:"subject"`
	Data    PaymentCreatedEventData `json:"data"`
}

func NewPaymentCreatedEvent(id string, orderId string, stripeId string) *PaymentCreatedEvent {
	return &PaymentCreatedEvent{
		Subject: PaymentCreated,
		Data: PaymentCreatedEventData{
			Id:       id,
			OrderId:  orderId,
			StripeId: stripeId,
		},
	}
}
