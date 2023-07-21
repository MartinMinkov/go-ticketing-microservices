package state

import "github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/database"

type AppState struct {
	DB        *database.Database
	DBCleanup func()
}
