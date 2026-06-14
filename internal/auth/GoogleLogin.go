package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GoogleUserLogin(c *gin.Context) {
	url := GoogleConfig.AuthCodeURL(
		"user",
		oauth2.AccessTypeOffline,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleOrganizationLogin(c *gin.Context) {
	url := GoogleConfig.AuthCodeURL(
		"organization",
		oauth2.AccessTypeOffline,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}