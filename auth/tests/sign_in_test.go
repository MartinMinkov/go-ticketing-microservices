package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/utils"
)

func TestSignInIsSuccessful(t *testing.T) {
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

	cookie, err := utils.FindCookie(resp)
	if err != nil {
		t.Fatal(err)
	}

	resp = app.POST_sign_in(data, cookie)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}

	var body map[string]interface{}
	if err := ReadJSON(resp, &body); err != nil {
		t.Fatal(err)
	}

	if body["email"] != data["email"] {
		t.Fatalf("Expected email %v, got %v", data["email"], body["email"])
	}
}

func TestSignInWithWrongEmail(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"email":    "test@test.com",
		"password": "pass",
	}

	resp := app.POST_sign_up(data)
	defer resp.Body.Close()

	cookie, err := utils.FindCookie(resp)
	if err != nil {
		t.Fatal(err)
	}

	data["email"] = "test1@gmail.com"
	resp = app.POST_sign_in(data, cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status code %v, got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestSignInWithWrongPassword(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"email":    "test@test.com",
		"password": "pass",
	}

	resp := app.POST_sign_up(data)
	defer resp.Body.Close()

	cookie, err := utils.FindCookie(resp)
	if err != nil {
		t.Fatal(err)
	}

	data["password"] = "wrong"
	resp = app.POST_sign_in(data, cookie)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status code %v, got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}
