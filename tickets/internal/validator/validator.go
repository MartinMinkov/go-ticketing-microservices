package validator

import (
	"fmt"

	v "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TicketInput interface {
	Title() string
	Price() int64
}

func validateInput(input TicketInput, validateFields map[string]string) ([]v.ErrorInnerResponse, error) {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errMessages := make([]v.ErrorInnerResponse, len(errs))
		for i, e := range errs {
			var errMsg string
			if validateMsg, ok := validateFields[e.Field()]; ok {
				errMsg = validateMsg
			} else {
				errMsg = "Unknown field error"
			}
			errMessages[i] = v.ErrorInnerResponse{Msg: errMsg, Field: e.Field(), Value: fmt.Sprint(e.Value())}
		}
		return errMessages, err
	}
	return nil, nil
}

func ValidateTicket(c *gin.Context, input TicketInput) error {
	validateFields := map[string]string{
		"ID":    "ID validation error",
		"Title": "Title validation error",
		"Price": "Price validation error",
	}
	errMessages, err := validateInput(input, validateFields)
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
	errMessages, err := validateInput(input, validateFields)
	if err != nil {
		errors := v.ErrorResponse{Errors: errMessages}
		errors.Error(c)
	}
	return err
}
