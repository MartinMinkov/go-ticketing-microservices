package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/auth/api"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/state"
	"github.com/gin-gonic/gin"
)

type TestApp struct {
	Config   *config.Config
	AppState *state.AppState
	Server   *http.Server
}

func (t *TestApp) GetHealthCheck() *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	resp, err := http.Get(addr + "/api/users/healthcheck")
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) GetCurrentUser(cookie *http.Cookie) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodGet, addr+"/api/users/currentuser", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) PostSignUp(data map[string]interface{}) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(addr+"/api/users/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) PostSignIn(data map[string]interface{}, cookie *http.Cookie) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, addr+"/api/users/signin", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) PostSignOut(cookie *http.Cookie) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodPost, addr+"/api/users/signout", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func SpawnApp() *TestApp {
	applicationConfig := config.BuildApplicationConfig()
	databaseConfig := config.BuildDatabaseConfig()
	databaseConfig.Database = "auth_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)

	config := &config.Config{
		ApplicationConfig: applicationConfig,
		DatabaseConfig:    databaseConfig,
	}

	appState := api.BuildAppState(config)

	gin.SetMode(gin.ReleaseMode)
	server := api.BuildServer(config, appState)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait a bit for the server to start
	time.Sleep(500 * time.Millisecond)

	return &TestApp{
		Config:   config,
		AppState: appState,
		Server:   server,
	}
}

func (app *TestApp) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer app.AppState.DBCleanup()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}
