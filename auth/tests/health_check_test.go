package tests

import (
	"net/http"
	"testing"
)

func TestHealthCheckReturns200(t *testing.T) {
	app := SpawnApp()
	defer app.Cleanup()

	resp := app.GET_healthcheck()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := ReadJSON(resp, &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "UP" {
		t.Fatalf("Expected status ok, got %v", body["status"])
	}
}
