package routes

import (
	"auth.mminkov.net/internal/state"
	"github.com/gin-gonic/gin"
)

func Healthcheck(c *gin.Context, appState *state.AppState) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
