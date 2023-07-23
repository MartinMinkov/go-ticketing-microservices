package routes

import (
	"context"
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CreateTicket(c *gin.Context, appState *state.AppState) {
	userClaims := middleware.GetUserClaimsByContext(c)

	var createticketInput createTicketInput
	if err := c.ShouldBindJSON(&createticketInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := validator.ValidateTicket(c, createticketInput)
	if err != nil {
		log.Err(err).Msg("Failed to validate ticket")
		return
	}

	existingTicket, _ := model.GetSingleTicketByTitle(appState.DB, createticketInput.Title())
	if existingTicket != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket already exists"})
		return
	}

	ticket := model.NewTicket(userClaims.ID, createticketInput.Title(), createticketInput.Price())
	ticket.Save(appState.DB)

	publisher := events.NewPublisher(appState.NatsConn, events.TicketCreated, context.TODO())
	err = publisher.Publish(events.NewTicketCreatedEvent(ticket.ID.Hex(), ticket.UserId, ticket.Title, ticket.Price))
	if err != nil {
		log.Err(err).Msg("Failed to publish ticket created event")
	}

	c.JSON(http.StatusCreated, ticket)
}

type createTicketInput struct {
	TitleField string `json:"title" validate:"required,min=1,max=100"`
	PriceField int64  `json:"price" validate:"required,gt=0"`
}

func (c createTicketInput) Title() string {
	return c.TitleField
}

func (c createTicketInput) Price() int64 {
	return c.PriceField
}
