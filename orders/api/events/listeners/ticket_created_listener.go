package listeners

import (
	"context"
	"encoding/json"
	"log"
	"time"

	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const QueueGroupName = "orders-service"

type TicketCreatedListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *TicketCreatedListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var ticketCreatedEvent e.TicketCreatedEvent
	err := json.Unmarshal(msg.Data(), &ticketCreatedEvent)
	if err != nil {
		log.Default().Println("listener: Could not unmarshal data", err)
		msg.Ack()
		return nil, err
	}
	return ticketCreatedEvent, nil
}

func (t *TicketCreatedListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	ticketCreatedEvent, ok := data.(e.TicketCreatedEvent)

	defer func() {
		msg.Ack()
	}()

	if !ok {
		log.Default().Println("listener: Could not cast data to TicketCreatedEvent")
		return nil
	}

	ticket := model.NewTicket(ticketCreatedEvent.Data.Id, ticketCreatedEvent.Data.Title, ticketCreatedEvent.Data.Price)
	ticket.Version = ticketCreatedEvent.Data.Version
	err := ticket.Save(t.db)
	if err != nil {
		log.Default().Println("listener: Could not save ticket in DB", err)
		return err
	}

	return nil
}

func NewTicketCreatedListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *TicketCreatedListener {
	t := &TicketCreatedListener{}
	listener := e.NewListener(ns, e.TicketCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
