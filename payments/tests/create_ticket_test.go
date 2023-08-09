package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
)

// func TestCreateTicketIsSuccessful(t *testing.T) {
// 	app := SpawnApp()
// 	defer app.Cleanup()

// 	data := map[string]interface{}{
// 		"title": "test ticket",
// 		"price": 100,
// 	}

// 	user := utils.UserMock1
// 	resp := app.PostCreateTicket(data, &user)
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		t.Errorf("expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
// 	}
// }

func TestCreatePaymentFailsWhenOrderIdIsMissing(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"token": "test token",
	}

	user := utils.UserMock1
	resp := app.PostCreatePayment(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateTicketFailsWhenTokenIsMissing(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"order_id": "test order id",
	}

	user := utils.UserMock1
	resp := app.PostCreatePayment(data, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateTicketFailsWhenUserIsUnauthenticated(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	user := utils.UserMock1
	user.JWT = ""
	resp := app.PostCreatePayment(nil, &user)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}
