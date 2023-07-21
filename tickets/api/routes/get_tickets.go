package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/gin-gonic/gin"
)

func GetTickets(c *gin.Context, appState *state.AppState) {
	userClaims := middleware.GetUserClaimsByContext(c)
	if userClaims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tickets, err := model.GetAllTickets(appState.DB, userClaims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}
