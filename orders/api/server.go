package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
	e "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/api/events/listeners"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/api/routes"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func BuildServer(cfg *config.Config, appState *state.AppState) *http.Server {
	InitLogger()
	r := initGin(appState)

	if cfg.ApplicationConfig.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		log.Err(err).Msg("Failed to get JWT secret")
	}
	auth.SetSecret(jwtSecret)

	return &http.Server{
		Addr:    ":" + cfg.ApplicationConfig.Port,
		Handler: r,
	}
}

func InitEventListeners(appState *state.AppState) {
	ticketCreatedListener := listeners.NewTicketCreatedListener(appState.NatsConn, time.Second*5, appState.DB, context.TODO())
	ticketCreatedListener.Listener.Listen()

	// ticketUpdatedListener := listeners.NewTicketUpdatedListener(appState.NatsConn, time.Second*5, appState.DB, appState.DB.Ctx)
	// ticketUpdatedListener.Listener.Listen()
}

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func BuildAppState(config *config.Config) *state.AppState {
	db := database.ConnectDB(config)

	nc, err := e.ConnectWithRetry(config.NatsConfig.GetAddress(), time.Second*5, time.Second*15)
	if err != nil {
		panic("Failed to connect to NATS")
	}

	cleanup := func() {
		if err := db.Client.Disconnect(db.Ctx); err != nil {
			panic("Failed to disconnect from MongoDB")
		}
	}

	return &state.AppState{DB: db, DBCleanup: cleanup, NatsConn: nc, NatsCleanup: func() { nc.Close() }}
}

func initGin(appState *state.AppState) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger("orders"))
	setupRoutes := routes.Routes(appState)
	setupRoutes(r)

	return r
}