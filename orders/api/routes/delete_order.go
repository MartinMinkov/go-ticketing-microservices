package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"
	"github.com/gin-gonic/gin"
)

func DeleteOrder(c *gin.Context, appState *state.AppState) {
	orderId := c.Param("id")
	err := model.CancelOrder(appState.DB, orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
