package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
)

func TestCurrentUserRespondsWithCorrectDetails(t *testing.T) {
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

	var body map[string]interface{}
	if err := ReadJSON(resp, &body); err != nil {
		t.Fatal(err)
	}

	cookie, err := auth.FindCookie(resp)
	if err != nil {
		t.Fatal(err)
	}

	resp = app.GET_currentuser(cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}

	if err := ReadJSON(resp, &body); err != nil {
		t.Fatal(err)
	}

	if body["email"] == nil {
		t.Fatalf("Expected current user to be returned")
	}

	if body["email"] != data["email"] {
		t.Fatalf("Expected email to be %v, got %v", data["email"], body["email"])
	}
}

func TestCurrentUserFailsWithBadCookie(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	cookie := &http.Cookie{
		Name:  auth.JWT_COOKIE_NAME,
		Value: "bad",
	}

	resp := app.GET_currentuser(cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status code %v, got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}
