package routes

import (
	"context"
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func DeleteOrder(c *gin.Context, appState *state.AppState) {
	orderId := c.Param("id")
	cancelledOrder, err := model.CancelOrder(appState.DB, orderId)
	if err != nil {
		log.Err(err).Msg("Failed to cancel order")
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	publisher := events.NewPublisher(appState.NatsConn, events.OrderCancelled, context.TODO())
	err = publisher.Publish(events.NewOrderCancelledEvent(orderId, *cancelledOrder.TicketId))
	if err != nil {
		log.Err(err).Msg("Failed to publish order cancelled event")
	}
	c.JSON(http.StatusNoContent, nil)
}
