package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GoogleLogin(c *gin.Context) {
	url := GoogleConfig.AuthCodeURL(
		"random-state",
		oauth2.AccessTypeOffline,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}