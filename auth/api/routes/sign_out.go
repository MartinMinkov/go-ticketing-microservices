package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
	"github.com/gin-gonic/gin"
)

func SignOut(c *gin.Context) {
	auth.DeleteCookieHandler(c)
	c.JSON(http.StatusOK, nil)
}
