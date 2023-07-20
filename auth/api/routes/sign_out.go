package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/utils"
	"github.com/gin-gonic/gin"
)

func SignOut(c *gin.Context) {
	utils.DeleteCookieHandler(c)
	c.JSON(http.StatusOK, nil)
}
