package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/utils"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/api"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/gin-gonic/gin"
)

type TestApp struct {
	Config   *config.Config
	AppState *state.AppState
	Server   *http.Server
}

func (t *TestApp) GetHealthCheck() *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	resp, err := http.Get(addr + "/api/tickets/healthcheck")
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) GetTicket(ticketId string, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodGet, addr+"/api/tickets/"+ticketId, nil)
	if err != nil {
		log.Fatalln(err)
	}

	utils.SetUserCookieOnRequest(req, user)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) GetTickets(user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodGet, addr+"/api/tickets", nil)
	if err != nil {
		log.Fatalln(err)
	}

	utils.SetUserCookieOnRequest(req, user)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) PostCreateTicket(data map[string]interface{}, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, addr+"/api/tickets", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	utils.SetUserCookieOnRequest(req, user)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) PutUpdateTicket(data map[string]interface{}, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPut, addr+"/api/tickets", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	utils.SetUserCookieOnRequest(req, user)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func (t *TestApp) DeleteTicket(ticketId string, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodDelete, addr+"/api/tickets/"+ticketId, nil)
	if err != nil {
		log.Fatalln(err)
	}

	utils.SetUserCookieOnRequest(req, user)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func SpawnApp() *TestApp {
	applicationConfig := config.BuildApplicationConfig()
	natsConfig := config.BuildNatsConfig()
	databaseConfig := config.BuildDatabaseConfig()

	databaseConfig.Database = "tickets_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	databaseConfig.Port = 27018 // When using docker compose, increment the expected port by 1

	config := &config.Config{
		ApplicationConfig: applicationConfig,
		DatabaseConfig:    databaseConfig,
		NatsConfig:        natsConfig,
	}

	appState := api.BuildAppState(config)
	MockTicketsInDatabase(appState)

	gin.SetMode(gin.ReleaseMode)
	server := api.BuildServer(config, appState)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait a bit for the server to start
	time.Sleep(100 * time.Millisecond)

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

func MockTicket(userId string, appState *state.AppState) {
	t := model.NewTicket(userId, "test", 1)
	err := t.Save(appState.DB)
	if err != nil {
		log.Fatal(err)
	}
}

func MockTicketsInDatabase(appState *state.AppState) {
	// Mock 5 tickets for user 1
	for i := 0; i < 5; i++ {
		MockTicketForUser1(appState)
	}

	// Mock 3 tickets for user 2
	for i := 0; i < 3; i++ {
		MockTicketForUser2(appState)
	}

	// Mock 1 ticket for user 3
	MockTicketForUser3(appState)
}

func MockTicketForUser1(appState *state.AppState) {
	user1 := utils.UserMock1
	MockTicket(user1.ID, appState)
}

func MockTicketForUser2(appState *state.AppState) {
	user2 := utils.UserMock1
	MockTicket(user2.ID, appState)
}

func MockTicketForUser3(appState *state.AppState) {
	user3 := utils.UserMock1
	MockTicket(user3.ID, appState)
}
