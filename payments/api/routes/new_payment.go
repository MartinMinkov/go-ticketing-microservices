package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/events"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/charge"
)

func CreatePayment(c *gin.Context, appState *state.AppState) {
	userClaims := middleware.GetUserClaimsByContext(c)

	var input createPaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Err(err).Msg("Failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := validator.ValidatePayment(c, input)
	if err != nil {
		log.Err(err).Msg("Failed to validate payment")
		return
	}

	existingOrder, err := model.GetSingleOrder(appState.DB, input.ID())
	if err != nil {
		log.Err(err).Msg("Failed to get order")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch order"})
		return
	}
	if existingOrder == nil {
		log.Err(err).Msg("Order does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Order does not exist"})
		return
	}
	if *existingOrder.UserId != userClaims.ID {
		log.Err(err).Msg("Order does not belong to user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Order does not belong to user"})
		return
	}
	if *existingOrder.Status == "cancelled" {
		log.Err(err).Msg("Order is cancelled")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is cancelled"})
		return
	}

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(*existingOrder.Price * 100),
		Currency:    stripe.String(string(stripe.CurrencyCAD)),
		Description: stripe.String(fmt.Sprintf("OrderID: %s", input.ID())),
		Source:      &stripe.PaymentSourceSourceParams{Token: stripe.String(input.Token())},
	}

	_, err = charge.New(params)
	if err != nil {
		log.Err(err).Msg("Failed to create charge")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create charge"})
		return
	}

	payment := model.NewPayment(existingOrder.ID.Hex(), input.Token())
	err = payment.Save(appState.DB)
	if err != nil {
		log.Err(err).Msg("Failed to save payment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment"})
		return
	}

	publisher := events.NewPublisher(appState.NatsConn, events.PaymentCreated, context.TODO())
	err = publisher.Publish(events.NewPaymentCreatedEvent(payment.ID.Hex(), payment.OrderId, payment.StripeId))
	if err != nil {
		log.Err(err).Msg("Failed to publish order created event")
	}

	c.JSON(http.StatusOK, gin.H{"id": payment.ID.Hex()})
}

type createPaymentInput struct {
	OrderId    string `json:"order_id" binding:"required"`
	TokenField string `json:"token" binding:"required"`
}

func (i createPaymentInput) ID() string {
	return i.OrderId
}

func (i createPaymentInput) Token() string {
	return i.TokenField
}
