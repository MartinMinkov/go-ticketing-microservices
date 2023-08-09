package state

import (
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/database"
	"github.com/nats-io/nats.go"
	"github.com/stripe/stripe-go/v74/client"
)

type AppState struct {
	DB           *database.Database
	DBCleanup    func()
	NatsConn     *nats.Conn
	NatsCleanup  func()
	StripeClient *client.API
}
