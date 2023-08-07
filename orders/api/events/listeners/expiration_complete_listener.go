package listeners

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type ExpirationCompleteListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *ExpirationCompleteListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var expirationCompleteEvent e.ExpirationCompleteEvent
	err := json.Unmarshal(msg.Data(), &expirationCompleteEvent)
	if err != nil {
		log.Default().Println("listener: Could not unmarshal data", err)
		return nil, err
	}
	return expirationCompleteEvent, nil
}

func (t *ExpirationCompleteListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	expirationCompleteEvent, ok := data.(e.ExpirationCompleteEvent)
	if !ok {
		log.Default().Println("listener: Could not cast data to ExpirationCompleteEvent")
		return nil
	}

	order, err := model.GetSingleOrder(t.db, expirationCompleteEvent.Data.OrderId)
	if err != nil {
		log.Default().Println("listener: Could not get order", err)
		return err
	}

	orderId := order.ID.Hex()
	orderCancelledStatus := string(model.OrderCancelled)
	order.Status = &orderCancelledStatus
	err = order.Update(t.db)
	if err != nil {
		log.Default().Println("listener: Could not set order to cancelled", err)
		return err
	}

	publisher := events.NewPublisher(t.Listener.Ns, events.OrderCancelled, context.TODO())
	err = publisher.Publish(events.NewOrderCancelledEvent(orderId, *order.TicketId, *order.Version))
	if err != nil {
		log.Default().Println("listener: Could not publish order cancelled event", err)
	}

	msg.Ack()
	log.Default().Println("listener: Acknowledged message", orderId)
	return nil
}

func NewExpirationCompleteListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *ExpirationCompleteListener {
	t := &ExpirationCompleteListener{}
	listener := e.NewListener(ns, e.ExpirationComplete, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
