package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
)

func TestGetTicketIsSuccessful(t *testing.T) {
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

	if ticket.ID.Hex() == "" {
		t.Error("expected ticket to be created")
	}

	resp = app.GetTicket(ticket.ID.Hex(), &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetTicketReturnsNotFoundWhenTicketDoesNotExist(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	resp := app.GetTicket("get_me_the_ticket", &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestGetTicketFailsWhenUserIsUnauthenticated(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	user.JWT = ""
	resp := app.GetTickets(&user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
