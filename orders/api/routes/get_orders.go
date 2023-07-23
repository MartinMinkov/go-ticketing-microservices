package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"
	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context, appState *state.AppState) {
	userClaims := middleware.GetUserClaimsByContext(c)
	orders, err := model.GetAllOrders(appState.DB, userClaims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
