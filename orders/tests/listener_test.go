package tests

import (
	"context"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTicketCreatedListenerCreatesTicket(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	// Publish the event with a random ticket id
	ticket_id := primitive.NewObjectID().Hex()
	publisher := events.NewPublisher(app.AppState.NatsConn, events.TicketCreated, context.TODO())
	err := publisher.Publish(events.NewTicketCreatedEvent(ticket_id, primitive.NewObjectID().Hex(), "Test Ticket", 100, 0))
	if err != nil {
		log.Err(err).Msg("Failed to publish ticket created event")
	}
	// Wait for the event to be processed
	app.Wait(100)

	ticket, err := model.GetSingleTicket(app.AppState.DB, ticket_id)
	if err != nil {
		t.Fatal(err)
	}
	if ticket.Title != "Test Ticket" {
		t.Errorf("expected title %s, got %s", "Test Ticket", ticket.Title)
	}
}

func TestTicketUpdatedListenerUpdatesTicket(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	// Publish the event with a random ticket id
	ticket_id := primitive.NewObjectID().Hex()
	publisher := events.NewPublisher(app.AppState.NatsConn, events.TicketCreated, context.TODO())
	err := publisher.Publish(events.NewTicketCreatedEvent(ticket_id, primitive.NewObjectID().Hex(), "Test Ticket", 100, 0))
	if err != nil {
		log.Err(err).Msg("Failed to publish ticket created event")
	}

	app.Wait(250)

	publisher1 := events.NewPublisher(app.AppState.NatsConn, events.TicketUpdated, context.TODO())
	err = publisher1.Publish(events.NewTicketUpdatedEvent(ticket_id, primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), "Test Ticket1", 110, 1))
	if err != nil {
		log.Err(err).Msg("Failed to publish ticket update event")
	}

	app.Wait(250)

	ticket, err := model.GetSingleTicket(app.AppState.DB, ticket_id)
	if err != nil {
		t.Fatal(err)
	}
	if ticket.Title != "Test Ticket1" {
		t.Errorf("expected title %s, got %s", "Test Ticket", ticket.Title)
	}
	if ticket.Price != 110 {
		t.Errorf("expected price %d, got %d", 110, ticket.Price)
	}
	if ticket.Version != 1 {
		t.Errorf("expected version %d, got %d", 1, ticket.Version)
	}
}
