package validator

import (
	v "github.com/MartinMinkov/go-ticketing-microservices/common/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateEmailPassword(c *gin.Context, input interface {
	Email() string
	Password() string
}) error {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errMessages := make([]v.ErrorInnerResponse, len(errs))
		for i, e := range errs {
			var errMsg string
			switch e.Field() {
			case "Email":
				errMsg = "Email validation error"
			case "Password":
				errMsg = "Password validation error"
			default:
				errMsg = "Unknown field error"
			}
			errMessages[i] = v.ErrorInnerResponse{Msg: errMsg, Field: e.Field(), Value: e.Value().(string)}
		}
		errors := v.ErrorResponse{Errors: errMessages}
		errors.Error(c)
	}
	return err
}
