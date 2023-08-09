package validator

import (
	v "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/validator"
	"github.com/gin-gonic/gin"
)

type OrderInput interface {
	ID() string
	Token() string
}

func ValidatePayment(c *gin.Context, input OrderInput) error {
	validateFields := map[string]string{
		"ID":    "ID validation error",
		"Token": "Token validation error",
	}
	errMessages, err := v.ValidateInput(input, validateFields)
	if err != nil {
		errors := v.ErrorResponse{Errors: errMessages}
		errors.Error(c)
	}
	return err
}
