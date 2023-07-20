package routes

import (
	"net/http"

	"auth.mminkov.net/internal/utils"
	"github.com/gin-gonic/gin"
)

func SignOut(c *gin.Context) {
	utils.DeleteCookieHandler(c)
	c.JSON(http.StatusOK, nil)
}
