package routes

import (
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/state"
	"github.com/gin-gonic/gin"
)

func Routes(appState *state.AppState) func(router *gin.Engine) {
	return func(router *gin.Engine) {
		router.GET("/api/tickets/healthcheck", func(c *gin.Context) {
			Healthcheck(c, appState)
		})

		authorized := router.Group("/api/tickets")
		{
			authorized.Use(middleware.JWTMiddleware())
			authorized.Use(middleware.UserMiddleware())

			authorized.GET("/", func(c *gin.Context) {
				GetTickets(c, appState)
			})
			authorized.GET("/:id", func(c *gin.Context) {
				GetTicket(c, appState)
			})

			authorized.POST("/", func(c *gin.Context) {
				CreateTicket(c, appState)
			})

			authorized.PUT("/", func(c *gin.Context) {
				UpdateTicket(c, appState)
			})

			authorized.DELETE("/:id", func(c *gin.Context) {
				DeleteTicket(c, appState)
			})
		}
	}
}
