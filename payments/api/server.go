package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/api/events/listeners"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/api/routes"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/config"
	c "github.com/MartinMinkov/go-ticketing-microservices/payments/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/state"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func BuildServer(cfg *config.Config, appState *state.AppState) *http.Server {
	InitLogger()
	r := initGin(appState)

	if cfg.ApplicationConfig.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtSecret, err := c.GetJWTSecret()
	if err != nil {
		log.Err(err).Msg("Failed to get JWT secret")
	}
	auth.SetSecret(jwtSecret)

	return &http.Server{
		Addr:    ":" + cfg.ApplicationConfig.Port,
		Handler: r,
	}
}

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func InitEventListeners(appState *state.AppState) {
	go func() {
		orderCreatedListener := listeners.NewOrderCreatedListener(appState.NatsConn, time.Second*5, appState.DB, context.TODO())
		orderCreatedListener.Listener.Listen()
	}()
	go func() {
		orderCancelledListener := listeners.NewOrderCancelledListener(appState.NatsConn, time.Second*5, appState.DB, context.TODO())
		orderCancelledListener.Listener.Listen()
	}()
}

func BuildAppState(config *config.Config) *state.AppState {
	db := database.ConnectDB(config)

	nc, err := events.ConnectWithRetry(config.NatsConfig.GetAddress(), time.Second*5, time.Second*60)
	if err != nil {
		panic("Failed to connect to NATS")
	}

	cleanup := func() {
		if err := db.Client.Disconnect(db.Ctx); err != nil {
			panic("Failed to disconnect from MongoDB")
		}
	}

	stripeSecret, err := c.GetStripeSecret()
	if err != nil {
		log.Err(err).Msg("Failed to get JWT secret")
	}
	client := client.New(stripeSecret, nil)
	// Set global key
	stripe.Key = stripeSecret

	return &state.AppState{DB: db, DBCleanup: cleanup, NatsConn: nc, NatsCleanup: func() { nc.Close() }, StripeClient: client}
}

func initGin(appState *state.AppState) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger("payments"))
	setupRoutes := routes.Routes(appState)
	setupRoutes(r)

	return r
}
