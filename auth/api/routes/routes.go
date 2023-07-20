package routes

import (
	"auth.mminkov.net/internal/state"
	"common.mminkov.net/pkg/middleware"
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

			router.POST("/api/users/signout", SignOut)
		}
	}
}
