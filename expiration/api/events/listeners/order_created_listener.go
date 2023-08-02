package listeners

import (
	"context"
	"encoding/json"
	"log"
	"time"

	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/expiration/internal/redis"
	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const QueueGroupName = "order-service"

type OrderCreatedListener struct {
	Listener    *e.Listener
	AsynqClient *asynq.Client
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

	log.Default().Println("listener: Received order created event in expiration service", orderCreatedEvent.Data.Id)

	redis.EnqueueCreateExpiration(t.AsynqClient, orderCreatedEvent.Data.Id, time.Until(orderCreatedEvent.Data.ExpiresAt))
	msg.Ack()
	return nil
}

func NewOrderCreatedListener(ns *nats.Conn, client *asynq.Client, ackWait time.Duration, ctx context.Context) *OrderCreatedListener {
	t := &OrderCreatedListener{}
	listener := e.NewListener(ns, e.OrderCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.AsynqClient = client
	return t
}
