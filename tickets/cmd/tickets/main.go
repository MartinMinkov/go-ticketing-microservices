package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/tickets/api"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/config"
)

func main() {
	config := config.BuildConfig()
	appState := api.BuildAppState(config)
	server := api.BuildServer(config, appState)
	defer func() {
		appState.DBCleanup()
		appState.NatsCleanup()
	}()

	api.InitEventListeners(appState)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
