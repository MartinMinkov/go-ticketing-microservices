package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderCreatedListenerReservesTicket(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"title": "test ticket",
		"price": 100,
	}

	user := utils.UserMock1
	resp := app.PostCreateTicket(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var ticket model.Ticket
	if err := utils.ReadJSON(resp, &ticket); err != nil {
		t.Fatal(err)
	}

	order_id := primitive.NewObjectID().Hex()
	publisher := events.NewPublisher(app.AppState.NatsConn, events.OrderCreated, context.TODO())
	err := publisher.Publish(events.NewOrderCreatedEvent(order_id, ticket.UserId, ticket.ID.Hex(), string(events.OrderCreated), 100, time.Now().Add(30*time.Minute)))
	if err != nil {
		log.Err(err).Msg("Failed to publish order created event")
	}

	app.Wait(100)
	reserved_ticket, err := model.GetSingleTicket(app.AppState.DB, ticket.ID.Hex())

	if err != nil {
		t.Fatal(err)
	}

	if *reserved_ticket.OrderId != order_id {
		t.Errorf("expected order id %s, got %s", order_id, *ticket.OrderId)
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

func TestTicketCreatedListenerDoesNotAckSkippedVersion(t *testing.T) {
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

	publisher2 := events.NewPublisher(app.AppState.NatsConn, events.TicketUpdated, context.TODO())
	err = publisher2.Publish(events.NewTicketUpdatedEvent(ticket_id, primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), "Test Ticket2", 120, 5))
	if err != nil {
		log.Err(err).Msg("Failed to publish ticket update event")
	}
	app.Wait(250)
	ticket, err := model.GetSingleTicket(app.AppState.DB, ticket_id)
	if err != nil {
		t.Fatal(err)
	}
	if ticket.Title != "Test Ticket1" {
		t.Errorf("expected title %s, got %s", "Test Ticket1", ticket.Title)
	}
	if ticket.Price != 110 {
		t.Errorf("expected price %d, got %d", 110, ticket.Price)
	}
	if ticket.Version != 1 {
		t.Errorf("expected version %d, got %d", 1, ticket.Version)
	}
}