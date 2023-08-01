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

func TestOrderUpdatesListenerUnReservesTicket(t *testing.T) {
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

	publisher = events.NewPublisher(app.AppState.NatsConn, events.OrderCancelled, context.TODO())
	err = publisher.Publish(events.NewOrderCancelledEvent(*reserved_ticket.OrderId, ticket.ID.Hex()))
	if err != nil {
		log.Err(err).Msg("Failed to publish order cancelled event")
	}

	app.Wait(100)

	unreserved_ticket, err := model.GetSingleTicket(app.AppState.DB, ticket.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	if unreserved_ticket.OrderId != nil {
		t.Errorf("expected order id %s, got %s", order_id, *ticket.OrderId)
	}
}
