package routes

import (
	"net/http"

	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/model"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/state"
	"github.com/MartinMinkov/go-ticketing-microservices/auth/internal/validator"
	"github.com/MartinMinkov/go-ticketing-microservices/common/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SignIn(c *gin.Context, appState *state.AppState) {
	var signinInput signInInput
	if err := c.ShouldBindJSON(&signinInput); err != nil {
		log.Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := validator.ValidateEmailPassword(c, signinInput)
	if err != nil {
		log.Err(err).Msg("Failed to validate email and password")
		return
	}

	existingUser, err := model.FindByEmail(appState.DB, signinInput.Email())
	if err != nil || existingUser == nil {
		log.Err(err).Msg("Failed to find user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	passwordsMatch := existingUser.ComparePassword(signinInput.Password())
	if !passwordsMatch {
		log.Err(err).Msg("Invalid email or password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	jwt, err := auth.CreateJWT(existingUser.ID.Hex(), *existingUser.Email)
	if err != nil {
		log.Err(err).Msg("Failed to create JWT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	auth.SetCookieHandler(c, jwt)

	c.JSON(http.StatusOK,
		existingUser,
	)
}

type signInInput struct {
	EmailField    string `json:"email" validate:"required,email"`
	PasswordField string `json:"password" validate:"required,min=4,max=20"`
}

func (s signInInput) Email() string {
	return s.EmailField
}

func (s signInInput) Password() string {
	return s.PasswordField
}
