package state

import (
	"github.com/hibiken/asynq"
	"github.com/nats-io/nats.go"
)

type AppState struct {
	AsynqClient *asynq.Client
	NatsConn    *nats.Conn
	NatsCleanup func()
}
