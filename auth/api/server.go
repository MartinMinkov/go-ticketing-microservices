package api

import (
	"net/http"
	"os"

	"github.com/MartinMinkov/go-ticketing-microservices/auth/api/routes"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/config"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/database"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
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

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func BuildAppState(config *config.Config) *state.AppState {
	db := database.ConnectDB(config)
	cleanup := func() {
		if err := db.Client.Disconnect(db.Ctx); err != nil {
			panic("Failed to disconnect from MongoDB")
		}
	}
	return &state.AppState{DB: db, DBCleanup: cleanup}
}

func initGin(appState *state.AppState) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger("auth"))
	setupRoutes := routes.Routes(appState)
	setupRoutes(r)

	return r
}
