package validator

import (
	v "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/validator"
	"github.com/gin-gonic/gin"
)

type TicketInput interface {
	Title() string
	Price() int64
}

type OrderInput interface {
	TicketId() string
}

func ValidateOrder(c *gin.Context, input OrderInput) error {
	validateFields := map[string]string{
		"TicketId": "TicketId validation error",
	}
	errMessages, err := v.ValidateInput(input, validateFields)
	if err != nil {
		errors := v.ErrorResponse{Errors: errMessages}
		errors.Error(c)
	}
	return err
}

func ValidateTicket(c *gin.Context, input TicketInput) error {
	validateFields := map[string]string{
		"ID":    "ID validation error",
		"Title": "Title validation error",
		"Price": "Price validation error",
	}
	errMessages, err := v.ValidateInput(input, validateFields)
	if err != nil {
		errors := v.ErrorResponse{Errors: errMessages}
		errors.Error(c)
	}
	return err
}

func ValidateUpdateTicket(c *gin.Context, input TicketInput) error {
	validateFields := map[string]string{
		"ID":      "ID validation error",
		"Title":   "Title validation error",
		"Price":   "Price validation error",
		"Version": "Version validation error",
	}
	errMessages, err := v.ValidateInput(input, validateFields)
	if err != nil {
		errors := v.ErrorResponse{Errors: errMessages}
		errors.Error(c)
	}
	return err
}
