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
	"github.com/MartinMinkov/go-ticketing-microservices/payments/api"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/state"
	"github.com/gin-gonic/gin"
)

type TestApp struct {
	Config   *config.Config
	AppState *state.AppState
	Server   *http.Server
}

func (t *TestApp) GetHealthCheck() *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	resp, err := http.Get(addr + "/api/payments/healthcheck")
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func (t *TestApp) PostCreatePayment(data map[string]interface{}, user *utils.UserMock) *http.Response {
	addr := t.Config.ApplicationConfig.GetAddress()
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodPost, addr+"/api/payments", bytes.NewBuffer(jsonData))
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

	databaseConfig.Database = "payments_test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	databaseConfig.Port = 27020 // When using docker compose, increment the expected port by 1

	config := &config.Config{
		ApplicationConfig: applicationConfig,
		DatabaseConfig:    databaseConfig,
		NatsConfig:        natsConfig,
	}

	appState := api.BuildAppState(config)

	gin.SetMode(gin.ReleaseMode)
	server := api.BuildServer(config, appState)

	api.InitEventListeners(appState)

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

func (app *TestApp) Wait(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (app *TestApp) Cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer app.AppState.DBCleanup()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
}
