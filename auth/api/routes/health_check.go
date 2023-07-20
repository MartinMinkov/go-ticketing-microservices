package routes

import (
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/state"
	"github.com/gin-gonic/gin"
)

func Healthcheck(c *gin.Context, appState *state.AppState) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
