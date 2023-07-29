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

type OrderCancelledListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *OrderCancelledListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var orderCancelledEvent e.OrderCancelledEvent
	err := json.Unmarshal(msg.Data(), &orderCancelledEvent)
	if err != nil {
		log.Default().Println("listener: Could not unmarshal data", err)
		msg.Ack()
		return nil, err
	}
	return orderCancelledEvent, nil
}

func (t *OrderCancelledListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	orderCancelledEvent, ok := data.(e.OrderCancelledEvent)

	defer func() {
		msg.Ack()
	}()

	if !ok {
		log.Default().Println("listener: Could not cast data to OrderCancelledEvent")
		return nil
	}

	ticket, err := model.GetSingleTicket(t.db, orderCancelledEvent.Data.Ticket.Id)
	if err != nil {
		log.Default().Println("listener: Could not get ticket from DB", err)
		return err
	}
	ticket.OrderId = nil
	err = ticket.Update(t.db)
	if err != nil {
		log.Default().Println("listener: Could not save ticket in DB", err)
		return err
	}

	publisher := events.NewPublisher(t.Listener.Ns, events.TicketUpdated, context.TODO())
	err = publisher.Publish(events.NewTicketUpdatedEvent(ticket.ID.Hex(), ticket.UserId, *ticket.OrderId, ticket.Title, ticket.Price, ticket.Version))
	if err != nil {
		log.Default().Println("listener: Could not publish ticket updated event", err)
		return err
	}

	return nil
}

func NewOrderCancelledListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *OrderCancelledListener {
	t := &OrderCancelledListener{}
	listener := e.NewListener(ns, e.OrderCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
