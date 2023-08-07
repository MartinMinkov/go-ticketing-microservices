package listeners

import (
	"context"
	"encoding/json"
	"log"
	"time"

	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/model"
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

	order := model.NewOrder(orderCreatedEvent.Data.Id, orderCreatedEvent.Data.UserId, orderCreatedEvent.Data.Status, orderCreatedEvent.Data.Ticket.Price, orderCreatedEvent.Data.Version)
	err := order.Update(t.db)
	if err != nil {
		log.Default().Println("listener: Could not save order in DB", err)
		return err
	}

	msg.Ack()
	return nil
}

func NewOrderCreatedListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *OrderCreatedListener {
	t := &OrderCreatedListener{}
	listener := e.NewListener(ns, e.OrderCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
