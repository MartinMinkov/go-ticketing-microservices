package tests

import (
	"net/http"
	"testing"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
)

func TestSignUpIsSuccessful(t *testing.T) {
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
}

func TestMultipleSignUpIsSuccessful(t *testing.T) {
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

	data = map[string]interface{}{
		"email":    "test1@test.com",
		"password": "pass",
	}

	resp = app.POST_sign_up(data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %v, got %v", http.StatusCreated, resp.StatusCode)
	}

	data = map[string]interface{}{
		"email":    "test2@test.com",
		"password": "pass",
	}

	resp = app.POST_sign_up(data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %v, got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestSignUpFailsWithInvalidEmail(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"email":    "@test.com",
		"password": "pass",
	}

	resp := app.POST_sign_up(data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %v, got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestSignUpFailsWithInvalidPassword(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	data := map[string]interface{}{
		"email":    "test@test.com",
		"password": "",
	}

	resp := app.POST_sign_up(data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %v, got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestSignUpFailsWithDuplicateSignUp(t *testing.T) {
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

	resp = app.POST_sign_up(data)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %v, got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestSetsCookieAfterSignUp(t *testing.T) {
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

	_, err := auth.FindCookie(resp)
	if err != nil {
		t.Fatalf("Expected to find JWT cookie, got %v", err)
	}
}
