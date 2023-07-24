package events

type Subjects string

const (
	TicketCreated  Subjects = "ticket:created"
	TicketUpdated  Subjects = "ticket:updated"
	OrderCreated   Subjects = "order:created"
	OrderCancelled Subjects = "order:cancelled"
)
