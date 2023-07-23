package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
)

func TestCreateOrderIsSuccessful(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"ticket_id": app.Mocks.TicketIds[0],
	}

	user := utils.UserMock1
	resp := app.PostCreateOrder(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestCreateOrderFailsWhenTicketIsReserved(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"ticket_id": app.Mocks.TicketIds[0],
	}

	user := utils.UserMock1
	resp := app.PostCreateOrder(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	resp = app.PostCreateOrder(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateTicketFailsWhenTicketIdIsMissing(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"ticket_id": "",
	}

	user := utils.UserMock1
	resp := app.PostCreateOrder(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateTicketFailsWhenUserIsUnauthenticated(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"ticket_id": app.Mocks.TicketIds[0],
	}

	user := utils.UserMock1
	user.JWT = ""
	resp := app.PostCreateOrder(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
