package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/gin-gonic/gin"
)

func Healthcheck(c *gin.Context, appState *state.AppState) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
