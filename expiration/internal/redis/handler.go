package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
)

type ExpirationHandler struct {
	NatsConn *nats.Conn
}

func (h *ExpirationHandler) HandleExpiration(ctx context.Context, t *asynq.Task) error {
	switch t.Type() {
	case ExpirationType:
		var p Payload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		log.Printf(" [*] Expire OrderID: %s", p.OrderId)

		publisher := events.NewPublisher(h.NatsConn, events.ExpirationCreated, context.TODO())
		err := publisher.Publish(events.NewExpirationCreatedEvent(p.OrderId))
		if err != nil {
			return fmt.Errorf("failed to publish expiration created event: %s", err)
		}
	default:
		return fmt.Errorf("unexpected task type: %s", t.Type())
	}
	return nil
}
