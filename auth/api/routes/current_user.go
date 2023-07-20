package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/model"
	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {
	user := model.GetUserByContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, user)
}
