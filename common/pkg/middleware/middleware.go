package middleware

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type UserClaims struct {
	ID    string  `json:"id"`
	Email string  `json:"email"`
	Iat   float64 `json:"iat"`
}

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("jwt")
		if !exists {
			log.Info().Msg("Failed to get JWT from context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT"})
			c.Abort()
			return
		}

		jwtToken, ok := claims.(*jwt.Token)
		if !ok {
			log.Info().Msg("Failed to convert JWT from context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT"})
			c.Abort()
			return
		}

		user := &UserClaims{
			ID:    jwtToken.Claims.(jwt.MapClaims)["id"].(string),
			Email: jwtToken.Claims.(jwt.MapClaims)["email"].(string),
			Iat:   jwtToken.Claims.(jwt.MapClaims)["iat"].(float64),
		}

		c.Set("user", user)
		c.Next()
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.GetCookieHandler(c)
		if err != nil {
			log.Err(err).Msg("JWT not found in cookie")
			c.JSON(http.StatusBadRequest, gin.H{"error": "JWT not found"})
			c.Abort()
			return
		}

		verifiedJwt, err := auth.VerifyJWT(token)
		if err != nil {
			log.Err(err).Msg("JWT not verified")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT"})
			c.Abort()
			return
		}

		// Store the verified JWT into the context for use in subsequent handlers
		c.Set("jwt", verifiedJwt)
		c.Next()
	}
}
