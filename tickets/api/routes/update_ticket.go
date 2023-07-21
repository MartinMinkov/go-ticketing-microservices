package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func UpdateTicket(c *gin.Context, appState *state.AppState) {
	userClaims := middleware.GetUserClaimsByContext(c)

	var updateticketInput updateTicketCreated
	if err := c.ShouldBindJSON(&updateticketInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Debug().Msgf("update ticket input: %+v", updateticketInput)

	err := validator.ValidateUpdateTicket(c, updateticketInput)
	if err != nil {
		log.Err(err).Msg("Failed to validate ticket")
		return
	}

	currentTicket, err := model.GetSingleTicket(appState.DB, updateticketInput.ID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ticket"})
		return
	}
	if currentTicket == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket does not exist"})
		return
	}

	if *currentTicket.UserId != userClaims.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized to update ticket"})
		return
	}

	currentTicket.Title = &updateticketInput.TitleField
	currentTicket.Price = &updateticketInput.PriceField
	currentTicket.Version = &updateticketInput.VersionField

	err = currentTicket.Update(appState.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, currentTicket)

}

type updateTicketCreated struct {
	IdField      string `json:"id" validate:"required"`
	TitleField   string `json:"title" validate:"required,min=1,max=100"`
	PriceField   int64  `json:"price"`
	VersionField int64  `json:"version"`
}

func (c updateTicketCreated) ID() string {
	return c.IdField
}

func (c updateTicketCreated) Title() string {
	return c.TitleField
}

func (c updateTicketCreated) Price() int64 {
	return c.PriceField
}

func (c updateTicketCreated) Version() int64 {
	return c.VersionField
}