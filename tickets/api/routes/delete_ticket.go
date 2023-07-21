package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/gin-gonic/gin"
)

func DeleteTicket(c *gin.Context, appState *state.AppState) {
	ticketId := c.Param("id")
	ticket, err := model.GetSingleTicket(appState.DB, ticketId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	ticket.Delete(appState.DB)
	c.JSON(http.StatusNoContent, nil)
}
