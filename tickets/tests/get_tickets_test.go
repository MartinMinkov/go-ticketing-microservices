package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
)

func TestGetAllTicketIsSuccessful(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	resp := app.GetTickets(&user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var tickets []model.Ticket
	if err := utils.ReadJSON(resp, &tickets); err != nil {
		t.Fatal(err)
	}

	if len(tickets) == 0 {
		t.Error("expected tickets to be returned")
	}
}

func TestGetTicketsFailsWhenUserIsUnauthenticated(t *testing.T) {
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
