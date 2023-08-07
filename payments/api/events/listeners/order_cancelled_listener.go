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

type OrderCancelledListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *OrderCancelledListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var orderCancelledEvent e.OrderCancelledEvent
	err := json.Unmarshal(msg.Data(), &orderCancelledEvent)
	if err != nil {
		log.Default().Println("listener: Could not unmarshal data", err)
		return nil, err
	}
	return orderCancelledEvent, nil
}

func (t *OrderCancelledListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	orderCancelledEvent, ok := data.(e.OrderCancelledEvent)
	if !ok {
		log.Default().Println("listener: Could not cast data to OrderCancelledEvent")
		return nil
	}

	order, err := model.GetSingleOrder(t.db, orderCancelledEvent.Data.Id)
	if err != nil {
		log.Default().Println("listener: Could not get order", err)
		return err
	}

	order, err = model.CancelOrder(t.db, order.ID.Hex())
	if err != nil {
		log.Default().Println("listener: Could not cancel order", err)
		return err
	}

	// publisher := events.NewPublisher(t.Listener.Ns, events.TicketUpdated, context.TODO())
	// err = publisher.Publish(events.NewTicketUpdatedEvent(order.ID.Hex(), order.UserId, *order.OrderId, order.Title, order.Price, order.Version))
	// if err != nil {
	// 	log.Default().Println("listener: Could not publish ticket updated event", err)
	// 	return err
	// }

	msg.Ack()
	return nil
}

func NewOrderCancelledListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *OrderCancelledListener {
	t := &OrderCancelledListener{}
	listener := e.NewListener(ns, e.OrderCancelled, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
