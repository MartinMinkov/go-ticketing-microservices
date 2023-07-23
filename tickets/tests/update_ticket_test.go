package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
)

func TestUpdateTicketIsSuccessful(t *testing.T) {
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
	data["title"] = "updated title"
	data["price"] = 200
	data["version"] = ticket.Version
	data["id"] = ticket.ID.Hex()

	resp = app.PutUpdateTicket(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestUpdateTicketFailsWhenTicketDoesNotExist(t *testing.T) {
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
	data["title"] = "updated title"
	data["price"] = 200
	data["version"] = ticket.Version
	data["id"] = ""

	resp = app.PutUpdateTicket(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestUpdateTicketFailsWhenUserIsUnauthenticated(t *testing.T) {
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
	data["title"] = "updated title"
	data["price"] = 200
	data["version"] = ticket.Version

	user.JWT = ""
	resp = app.PutUpdateTicket(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestUpdateTicketIncrementsVersion(t *testing.T) {
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
	data["title"] = "updated title"
	data["price"] = 200
	data["version"] = ticket.Version
	data["id"] = ticket.ID.Hex()

	resp = app.PutUpdateTicket(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var updatedTicket model.Ticket
	if err := utils.ReadJSON(resp, &updatedTicket); err != nil {
		t.Fatal(err)
	}

	if updatedTicket.Version != ticket.Version+1 {
		t.Errorf("expected version %d, got %d", ticket.Version+1, updatedTicket.Version)
	}
}

func TestUpdateTicketFailsWhenUserDoesNotOwnTicket(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"title": "test ticket",
		"price": 100,
	}

	user1 := utils.UserMock1
	resp := app.PostCreateTicket(data, &user1)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var ticket model.Ticket
	if err := utils.ReadJSON(resp, &ticket); err != nil {
		t.Fatal(err)
	}
	data["title"] = "updated title"
	data["price"] = 200
	data["version"] = ticket.Version
	data["id"] = ticket.ID.Hex()

	user2 := utils.UserMock2
	resp = app.PutUpdateTicket(data, &user2)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestUpdateTicketImplementsOCC(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"title": "test ticket",
		"price": 100,
	}

	user := utils.UserMock1
	resp := app.PostCreateTicket(data, &user)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	ticket := model.Ticket{}
	if err := utils.ReadJSON(resp, &ticket); err != nil {
		t.Fatal(err)
	}

	resp = app.GetTicket(ticket.ID.Hex(), &user)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var ticket1 model.Ticket
	if err := utils.ReadJSON(resp, &ticket1); err != nil {
		t.Fatal(err)
	}

	resp = app.GetTicket(ticket.ID.Hex(), &user)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var ticket2 model.Ticket
	if err := utils.ReadJSON(resp, &ticket2); err != nil {
		t.Fatal(err)
	}

	data1 := map[string]interface{}{
		"title":   "test ticket1",
		"price":   200,
		"id":      ticket1.ID.Hex(),
		"version": ticket1.Version,
	}

	resp = app.PutUpdateTicket(data1, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	data2 := map[string]interface{}{
		"title":   "test ticket2",
		"price":   300,
		"id":      ticket2.ID.Hex(),
		"version": ticket2.Version,
	}

	resp = app.PutUpdateTicket(data2, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
