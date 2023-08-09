package routes

import (
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/state"
	"github.com/gin-gonic/gin"
)

func Routes(appState *state.AppState) func(router *gin.Engine) {
	return func(router *gin.Engine) {
		router.GET("/api/payments/healthcheck", func(c *gin.Context) {
			Healthcheck(c, appState)
		})

		authorized := router.Group("/api/payments")
		{
			authorized.Use(middleware.JWTMiddleware())
			authorized.Use(middleware.UserMiddleware())

			authorized.POST("/", func(c *gin.Context) {
				CreatePayment(c, appState)
			})
		}
	}
}
