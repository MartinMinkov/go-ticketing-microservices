package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
)

func TestDeleteOrderIsSuccessful(t *testing.T) {
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

	resp = app.DeleteOrder(order.ID.Hex(), &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestDeleteOrderFailsWhenOrderDoesNotExist(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1

	resp := app.DeleteOrder("orderId", &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestDeleteOrderFailsWhenUserIsUnauthenticated(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"ticket_id": app.Mocks.TicketIds[0],
	}

	user := utils.UserMock1
	resp := app.PostCreateOrder(data, &user)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}
