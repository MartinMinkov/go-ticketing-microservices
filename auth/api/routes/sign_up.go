package routes

import (
	"net/http"

	"auth.mminkov.net/internal/model"
	"auth.mminkov.net/internal/state"
	"auth.mminkov.net/internal/utils"
	"auth.mminkov.net/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func SignUp(c *gin.Context, appState *state.AppState) {
	var signupInput signUpInput
	if err := c.ShouldBindJSON(&signupInput); err != nil {
		log.Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validator.ValidateEmailPassword(c, signupInput)

	existingUser, err := model.FindByEmail(appState.DB, signupInput.Email())
	if err != nil {
		log.Err(err).Msg("Failed to find user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if existingUser != nil {
		log.Err(err).Msg("User already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := utils.GenerateFromPassword(signupInput.Password())
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := model.NewUser(signupInput.Email(), hashedPassword)
	err = user.Save(appState.DB)
	if err != nil {
		log.Err(err).Msg("Failed to save user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := utils.CreateJWT(user.ID.Hex(), *user.Email)
	if err != nil {
		log.Err(err).Msg("Failed to create JWT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.SetCookieHandler(c, jwt)
	c.JSON(http.StatusCreated, user)
}

type signUpInput struct {
	EmailField    string `json:"email" validate:"required,email"`
	PasswordField string `json:"password" validate:"required,min=4,max=20"`
}

func (s signUpInput) Email() string {
	return s.EmailField
}

func (s signUpInput) Password() string {
	return s.PasswordField
}
