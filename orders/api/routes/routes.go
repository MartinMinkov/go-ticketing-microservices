package routes

import (
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/state"
	"github.com/gin-gonic/gin"
)

func Routes(appState *state.AppState) func(router *gin.Engine) {
	return func(router *gin.Engine) {
		router.GET("/api/orders/healthcheck", func(c *gin.Context) {
			Healthcheck(c, appState)
		})

		authorized := router.Group("/api/orders")
		{
			authorized.Use(middleware.JWTMiddleware())
			authorized.Use(middleware.UserMiddleware())

			authorized.GET("/", func(c *gin.Context) {
				GetOrders(c, appState)
			})
			authorized.GET("/:id", func(c *gin.Context) {
				GetOrder(c, appState)
			})

			authorized.POST("/", func(c *gin.Context) {
				CreateOrder(c, appState)
			})

			authorized.DELETE("/:id", func(c *gin.Context) {
				DeleteOrder(c, appState)
			})
		}
	}
}
