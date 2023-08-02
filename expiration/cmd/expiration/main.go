package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/MartinMinkov/go-ticketing-microservices/expiration/api"
	"github.com/MartinMinkov/go-ticketing-microservices/expiration/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/expiration/internal/redis"
	"github.com/hibiken/asynq"
)

func main() {
	config := config.BuildConfig()
	appState := api.BuildAppState(config)
	server := api.BuildServer(config)

	api.InitLogger()
	api.InitEventListeners(appState)

	defer func() {
		appState.NatsCleanup()
	}()

	handler := &redis.ExpirationHandler{
		NatsConn: appState.NatsConn,
	}

	mux := asynq.NewServeMux()
	mux.HandleFunc(redis.ExpirationType, handler.HandleExpiration)
	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
