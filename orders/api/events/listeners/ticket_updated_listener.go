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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketUpdatedListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *TicketUpdatedListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var ticketUpdatedEvent e.TicketUpdatedEvent
	err := json.Unmarshal(msg.Data(), &ticketUpdatedEvent)
	if err != nil {
		log.Default().Println("update listener: Could not unmarshal data", err)
		return nil, err
	}
	return ticketUpdatedEvent, nil
}

func (t *TicketUpdatedListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	ticketUpdatedEvent, ok := data.(e.TicketUpdatedEvent)

	if !ok {
		log.Default().Println("update listener: Could not cast data to []byte", data)
		return nil
	}

	ticket := model.NewTicket(ticketUpdatedEvent.Data.UserId, ticketUpdatedEvent.Data.Title, ticketUpdatedEvent.Data.Price)
	id, err := primitive.ObjectIDFromHex(ticketUpdatedEvent.Data.Id)
	if err != nil {
		log.Default().Println("update listener: Could not parse ticket id", err)
		return err
	}

	ticket.ID = id
	ticket.Version = ticketUpdatedEvent.Data.Version
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
