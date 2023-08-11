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

type PaymentCreatedListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *PaymentCreatedListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	var paymentCreatedEvent e.PaymentCreatedEvent
	err := json.Unmarshal(msg.Data(), &paymentCreatedEvent)
	if err != nil {
		log.Default().Println("update listener: Could not unmarshal data", err)
		return nil, err
	}
	return paymentCreatedEvent, nil
}

func (t *PaymentCreatedListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	paymentCreatedEvent, ok := data.(e.PaymentCreatedEvent)

	if !ok {
		log.Default().Println("update listener: Could not cast data to []byte", data)
		return nil
	}

	order, err := model.GetSingleOrder(t.db, paymentCreatedEvent.Data.OrderId)
	if err != nil {
		return err
	}
	if *order.Status == string(model.OrderComplete) {
		return nil
	}

	_, err = model.CompleteOrder(t.db, order.ID.Hex())
	if err != nil {
		return err
	}

	msg.Ack()
	return nil
}

func NewPaymentCreatedListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *PaymentCreatedListener {
	t := &PaymentCreatedListener{}
	listener := e.NewListener(ns, e.PaymentCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
