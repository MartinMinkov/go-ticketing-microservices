package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
)

func TestGetAllTicketIsSuccessful(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	resp := app.GetOrders(&user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var orders []model.Order
	if err := utils.ReadJSON(resp, &orders); err != nil {
		t.Fatal(err)
	}

	if len(orders) == 0 {
		t.Error("expected orders to be returned")
	}
}

func TestGetOrdersFailsWhenUserIsUnauthenticated(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	user.JWT = ""
	resp := app.GetOrders(&user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
