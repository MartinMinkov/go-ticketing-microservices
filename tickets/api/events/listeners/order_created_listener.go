package listeners

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const QueueGroupName = "tickets-service"

type OrderCreatedListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *OrderCreatedListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var orderCreatedEvent e.OrderCreatedEvent
	err := json.Unmarshal(msg.Data(), &orderCreatedEvent)
	if err != nil {
		log.Default().Println("listener: Could not unmarshal data", err)
		return nil, err
	}
	return orderCreatedEvent, nil
}

func (t *OrderCreatedListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	orderCreatedEvent, ok := data.(e.OrderCreatedEvent)
	if !ok {
		log.Default().Println("listener: Could not cast data to OrderCreatedEvent")
		return nil
	}

	ticket, err := model.GetSingleTicket(t.db, orderCreatedEvent.Data.Ticket.Id)
	if err != nil {
		log.Default().Println("listener: Could not get ticket from DB", err)
		return err
	}

	ticket.OrderId = &orderCreatedEvent.Data.Id
	err = ticket.Update(t.db)
	if err != nil {
		log.Default().Println("listener: Could not save ticket in DB", err)
		return err
	}

	msg.Ack()

	publisher := events.NewPublisher(t.Listener.Ns, events.TicketUpdated, context.TODO())
	err = publisher.Publish(events.NewTicketUpdatedEvent(ticket.ID.Hex(), ticket.UserId, *ticket.OrderId, ticket.Title, ticket.Price, ticket.Version))
	if err != nil {
		log.Default().Println("listener: Could not publish ticket updated event", err)
		return err
	}

	return nil
}

func NewOrderCreatedListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *OrderCreatedListener {
	t := &OrderCreatedListener{}
	listener := e.NewListener(ns, e.OrderCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
