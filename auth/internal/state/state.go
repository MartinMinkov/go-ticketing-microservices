package state

import "github.com/MartinMinkov/go-ticketing-microservices/auth/internal/database"

type AppState struct {
	DB        *database.Database
	DBCleanup func()
}
