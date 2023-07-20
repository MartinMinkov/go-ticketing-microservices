package utils

import (
	"net/http"
	"time"

	"auth.mminkov.net/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_COOKIE_NAME = "auth-service-jwt"

func CreateJWT(userId string, email string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userId,
		"email": email,
		"iat":   time.Now().UTC().Unix(),
	})
	return signJWT(t)
}

func signJWT(t *jwt.Token) (string, error) {
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		return "", err
	}

	s, err := t.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return s, nil
}

func VerifyJWT(j string) (*jwt.Token, error) {
	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		return nil, err
	}

	return jwt.Parse(j, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

func GetCookieHandler(c *gin.Context) (string, error) {
	cookie, err := c.Cookie(JWT_COOKIE_NAME)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

func FindCookie(resp *http.Response) (*http.Cookie, error) {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == JWT_COOKIE_NAME {
			return cookie, nil
		}
	}
	return nil, nil

}

func SetCookieHandler(c *gin.Context, jwt string) {
	c.SetCookie(JWT_COOKIE_NAME, jwt, (3600 * 12), "/", "", false, true)
}

func DeleteCookieHandler(c *gin.Context) {
	c.SetCookie(JWT_COOKIE_NAME, "", -1, "/", "", false, true)
}
