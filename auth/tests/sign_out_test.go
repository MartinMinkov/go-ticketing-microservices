package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
)

func TestSignOutIsSuccessful(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"email":    "test@test.com",
		"password": "pass",
	}

	resp := app.POST_sign_up(data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %v, got %v", http.StatusCreated, resp.StatusCode)
	}

	cookie, err := auth.FindCookie(resp)
	if err != nil {
		t.Fatal(err)
	}

	resp = app.POST_sign_out(cookie)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}

	cookie, err = auth.FindCookie(resp)
	if err != nil || cookie.Value != "" {
		t.Fatalf("Expected cookie to be cleared, got %v", cookie)
	}
}
