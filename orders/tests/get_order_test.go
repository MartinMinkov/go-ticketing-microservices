package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
)

func TestGetTicketIsSuccessful(t *testing.T) {
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

	var order model.Order
	if err := utils.ReadJSON(resp, &order); err != nil {
		t.Fatal(err)
	}

	if order.ID.Hex() == "" {
		t.Error("expected ticket to be created")
	}

	resp = app.GetOrder(order.ID.Hex(), &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetOrderReturnsNotFoundWhenOrderDoesNotExist(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	resp := app.GetOrder("get_me_the_ticket", &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestGetOrderFailsWhenUserIsUnauthenticated(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	user.JWT = ""
	resp := app.GetOrder("", &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
