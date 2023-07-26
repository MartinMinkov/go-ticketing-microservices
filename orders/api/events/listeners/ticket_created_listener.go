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

const QueueGroupName = "orders-service"

type TicketCreatedListener struct {
	Listener *e.Listener
	db       *database.Database
}

func (t *TicketCreatedListener) ParseMessage(msg jetstream.Msg) (interface{}, error) {
	log.Default().Println("ParseMessage")

	var ticketCreatedEvent e.TicketCreatedEvent

	log.Default().Println("Raw Message Data: ", string(msg.Data()))

	err := json.Unmarshal(msg.Data(), &ticketCreatedEvent)
	if err != nil {
		log.Default().Println("Ticket created event error1: ", err)
		return nil, err
	}
	log.Default().Println("Ticket Event Data: ", ticketCreatedEvent.Data.Id, ticketCreatedEvent.Data.Title, ticketCreatedEvent.Data.Price)
	return ticketCreatedEvent, nil
}

func (t *TicketCreatedListener) OnMessage(data interface{}, msg jetstream.Msg) error {
	log.Default().Println("OnMessage")
	jsonData, ok := data.([]byte)
	if !ok {
		return nil
	}

	log.Default().Println("Ticket created event received")
	log.Default().Println("Raw Message Data: ", string(jsonData))

	var ticketCreatedEvent e.TicketCreatedEvent
	err := json.Unmarshal(jsonData, &ticketCreatedEvent)
	if err != nil {
		log.Default().Println("Ticket created event error2: ", err)
		return err
	}

	log.Default().Println("Ticket Data: ", ticketCreatedEvent.Data.Id, ticketCreatedEvent.Data.Title, ticketCreatedEvent.Data.Price)

	ticket := model.NewTicket(ticketCreatedEvent.Data.Id, ticketCreatedEvent.Data.Title, ticketCreatedEvent.Data.Price)
	ticket.Save(t.db)

	msg.Ack()

	log.Default().Println("Ticket created event processed")

	return nil
}

func NewTicketCreatedListener(ns *nats.Conn, ackWait time.Duration, db *database.Database, ctx context.Context) *TicketCreatedListener {
	t := &TicketCreatedListener{}
	listener := e.NewListener(ns, e.TicketCreated, QueueGroupName, ackWait, t, t, ctx)
	t.Listener = listener
	t.db = db
	return t
}
