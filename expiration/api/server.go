package api

import (
	"context"
	"os"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/expiration/api/events/listeners"
	"github.com/MartinMinkov/go-ticketing-microservices/expiration/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/expiration/internal/state"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func BuildServer(config *config.Config) *asynq.Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.RedisConfig.GetAddress()},
		asynq.Config{Concurrency: 10},
	)
	if srv == nil {
		panic("Failed to create asynq server")
	}
	return srv
}

func InitEventListeners(appState *state.AppState) {
	go func() {
		orderCreatedListener := listeners.NewOrderCreatedListener(appState.NatsConn, appState.AsynqClient, time.Second*5, context.TODO())
		orderCreatedListener.Listener.Listen()
	}()
}

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func BuildAppState(config *config.Config) *state.AppState {
	nc, err := events.ConnectWithRetry(config.NatsConfig.GetAddress(), time.Second*5, time.Second*60)
	if err != nil {
		panic("Failed to connect to NATS")
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisConfig.GetAddress()})
	if client == nil {
		panic("Failed to create redis client")
	}
	return &state.AppState{NatsConn: nc, NatsCleanup: func() { nc.Close() }, AsynqClient: client}
}
