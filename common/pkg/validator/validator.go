package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
