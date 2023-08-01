package routes

import (
	"context"
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CreateOrder(c *gin.Context, appState *state.AppState) {
	userClaims := middleware.GetUserClaimsByContext(c)

	var createOrderInput createOrderInput
	if err := c.ShouldBindJSON(&createOrderInput); err != nil {
		log.Err(err).Msg("Failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := validator.ValidateOrder(c, createOrderInput)
	if err != nil {
		log.Err(err).Msg("Failed to validate order")
		return
	}

	existingTicket, err := model.GetSingleTicket(appState.DB, createOrderInput.TicketId())
	if err != nil {
		log.Err(err).Msg("Failed to get ticket")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch ticket"})
		return
	}

	isReserved, err := model.IsTicketReserved(appState.DB, existingTicket.ID.Hex())
	if err != nil {
		log.Err(err).Msg("Failed to check if ticket is reserved")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch ticket"})
		return
	}
	if isReserved {
		log.Err(err).Msg("Ticket is already reserved")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket is already reserved"})
		return
	}

	order := model.NewOrder(userClaims.ID, createOrderInput.TicketId())
	order.Save(appState.DB)

	publisher := events.NewPublisher(appState.NatsConn, events.OrderCreated, context.TODO())
	err = publisher.Publish(events.NewOrderCreatedEvent(order.ID.Hex(), *order.UserId, *order.TicketId, *order.Status, existingTicket.Price, *order.ExpiresAt))
	if err != nil {
		log.Err(err).Msg("Failed to publish order created event")
	}

	c.JSON(http.StatusCreated, order)
}

type createOrderInput struct {
	TicketIdField string `json:"ticket_id" validate:"required"`
}

func (c createOrderInput) TicketId() string {
	return c.TicketIdField
}
