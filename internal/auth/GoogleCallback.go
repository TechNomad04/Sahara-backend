package auth

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func(h *Handler) GoogleCallback(c *gin.Context) {

	code := c.Query("code")

	token, err := GoogleConfig.Exchange(
		context.Background(),
		code,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	client := GoogleConfig.Client(
		context.Background(),
		token,
	)

	resp, err := client.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo",
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	var user struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	json.NewDecoder(resp.Body).Decode(&user)

	


	myTokens, err := IssueToken(user.ID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = Persist(context.Background(), h.Redis, myTokens)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	SendTokens(c, myTokens)
}