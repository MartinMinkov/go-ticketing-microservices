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
	"github.com/MartinMinkov/go-ticketing-microservices/orders/api"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"

	"github.com/gin-gonic/gin"
)

type Mocks struct {
	TicketIds []string
}

type TestApp struct {
	Config   *config.Config
	AppState *state.AppState
	Server   *http.Server
	Mocks    *Mocks
}

func (t *TestApp) GetHealthCheck() *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	resp, err := http.Get(addr + "/api/orders/healthcheck")
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) GetOrder(orderId string, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodGet, addr+"/api/orders/"+orderId, nil)
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

func (t *TestApp) GetOrders(user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodGet, addr+"/api/orders", nil)
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

func (t *TestApp) PostCreateOrder(data map[string]interface{}, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, addr+"/api/orders", bytes.NewBuffer(jsonData))
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

func (t *TestApp) DeleteOrder(ticketId string, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	req, err := http.NewRequest(http.MethodDelete, addr+"/api/orders/"+ticketId, nil)
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

	databaseConfig.Database = "orders_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	databaseConfig.Port = 27019 // When using docker compose, increment the expected port by 1

	config := &config.Config{
		ApplicationConfig: applicationConfig,
		DatabaseConfig:    databaseConfig,
		NatsConfig:        natsConfig,
	}

	appState := api.BuildAppState(config)
	ticketIds := MockTicketsInDatabase(appState)

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
		Mocks: &Mocks{
			TicketIds: ticketIds,
		},
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

func MockTicket(userId string, appState *state.AppState) string {
	t := model.NewTicket(userId, "test", 1)
	err := t.Save(appState.DB)
	if err != nil {
		log.Fatal(err)
	}
	return t.ID.Hex()
}

func MockTicketsInDatabase(appState *state.AppState) []string {
	ticketIds := []string{}
	// Mock 5 tickets for user 1
	for i := 0; i < 5; i++ {
		id := MockTicketForUser1(appState)
		ticketIds = append(ticketIds, id)
	}

	// Mock 3 tickets for user 2
	for i := 0; i < 3; i++ {
		id := MockTicketForUser2(appState)
		ticketIds = append(ticketIds, id)
	}

	// Mock 1 ticket for user 3
	id := MockTicketForUser3(appState)
	ticketIds = append(ticketIds, id)

	return ticketIds
}

func MockTicketForUser1(appState *state.AppState) string {
	user1 := utils.UserMock1
	return MockTicket(user1.ID, appState)
}

func MockTicketForUser2(appState *state.AppState) string {
	user2 := utils.UserMock1
	return MockTicket(user2.ID, appState)
}

func MockTicketForUser3(appState *state.AppState) string {
	user3 := utils.UserMock1
	return MockTicket(user3.ID, appState)
}
