package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"
	"github.com/gin-gonic/gin"
)

func GetOrder(c *gin.Context, appState *state.AppState) {
	orderId := c.Param("id")

	ticket, err := model.GetSingleOrder(appState.DB, orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}
