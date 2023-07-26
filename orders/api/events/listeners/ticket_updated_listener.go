package listeners

import (
	"context"
	"encoding/json"
	"time"

	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type TicketUpdatedListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *TicketUpdatedListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var ticketCreatedEvent e.TicketUpdatedEvent
	err := json.Unmarshal(msg.Data(), &ticketCreatedEvent)
	if err != nil {
		return nil, err
	}
	return ticketCreatedEvent, nil
}

func (t *TicketUpdatedListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	jsonData, ok := data.([]byte)
	if !ok {
		return nil
	}

	var ticketUpdatedEvent e.TicketUpdatedEvent
	err := json.Unmarshal([]byte(jsonData), &ticketUpdatedEvent)
	if err != nil {
		return err
	}

	ticket := model.NewTicket(ticketUpdatedEvent.Data.Id, ticketUpdatedEvent.Data.Title, ticketUpdatedEvent.Data.Price)
	err = ticket.Update(t.db)
	if err != nil {
		return err
	}

	msg.Ack()

	return nil
}

func NewTicketUpdatedListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *TicketUpdatedListener {
	t := &TicketUpdatedListener{}
	listener := e.NewListener(ns, e.TicketUpdated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
