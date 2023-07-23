package state

import (
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/database"
	"github.com/nats-io/nats.go"
)

type AppState struct {
	DB          *database.Database
	DBCleanup   func()
	NatsConn    *nats.Conn
	NatsCleanup func()
}
