package state

import "auth.mminkov.net/internal/database"

type AppState struct {
	DB        *database.Database
	DBCleanup func()
}
