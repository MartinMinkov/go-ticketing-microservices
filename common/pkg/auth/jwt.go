package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWT_COOKIE_NAME = "auth-service-jwt"

var JWT_SECRET = []byte("")

func SetSecret(secret string) {
	JWT_SECRET = []byte(secret)
}

func CreateJWT(userId string, email string, secret string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userId,
		"email": email,
		"iat":   time.Now().UTC().Unix(),
	})
	return signJWT(t)
}

func signJWT(t *jwt.Token) (string, error) {
	s, err := t.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}
	return s, nil
}

func VerifyJWT(j string) (*jwt.Token, error) {
	return jwt.Parse(j, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
}
