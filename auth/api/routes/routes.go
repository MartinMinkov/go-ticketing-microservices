package routes

import (
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(appState *state.AppState) func(router *gin.Engine) {
	return func(router *gin.Engine) {
		router.GET("/api/users/healthcheck", func(c *gin.Context) {
			Healthcheck(c, appState)
		})

		router.POST("/api/users/signup", func(c *gin.Context) {
			SignUp(c, appState)
		})

		router.POST("/api/users/signin", func(c *gin.Context) {
			SignIn(c, appState)
		})

		authorized := router.Group("/api/users")
		{
			authorized.Use(middleware.JWTMiddleware())
			authorized.Use(middleware.UserMiddleware())

			authorized.GET("/currentuser", CurrentUser)

			authorized.POST("/signout", SignOut)
		}
	}
}
