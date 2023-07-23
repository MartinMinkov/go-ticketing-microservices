package validator

import (
	"fmt"
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

func ValidateInput(input interface{}, validateFields map[string]string) ([]ErrorInnerResponse, error) {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		errMessages := make([]ErrorInnerResponse, len(errs))
		for i, e := range errs {
			var errMsg string
			if validateMsg, ok := validateFields[e.Field()]; ok {
				errMsg = validateMsg
			} else {
				errMsg = "Unknown field error"
			}
			errMessages[i] = ErrorInnerResponse{Msg: errMsg, Field: e.Field(), Value: fmt.Sprint(e.Value())}
		}
		return errMessages, err
	}
	return nil, nil
}
