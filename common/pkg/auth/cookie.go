package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
