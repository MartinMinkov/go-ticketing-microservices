package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors []ErrorInnerResponse `json:"errors"`
}

type ErrorInnerResponse struct {
	Msg   string `json:"msg"`
	Field string `json:"field"`
	Value string `json:"value"`
}

func (e *ErrorResponse) Error(c *gin.Context) {
	c.JSON(http.StatusBadRequest, e)
}

func ValidateEmailPassword(c *gin.Context, input interface {
	Email() string
	Password() string
}) {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errMessages := make([]ErrorInnerResponse, len(errs))
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
			errMessages[i] = ErrorInnerResponse{Msg: errMsg, Field: e.Field(), Value: e.Value().(string)}
		}
		errors := ErrorResponse{Errors: errMessages}
		errors.Error(c)
		return
	}
}
