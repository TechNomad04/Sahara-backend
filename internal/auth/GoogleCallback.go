package auth

import (
	"context"
	"encoding/json"
	"sahara/models"
	"strconv"

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
		GoogleId string `json:"id"`
		Name string `json:"name"`
		Email string `json:"email"`
	}

	json.NewDecoder(resp.Body).Decode(&user)

	var dbuser models.User

	if err := h.DB.
    Where("google_id = ?", user.GoogleId).
    First(&dbuser).Error; err != nil {
		dbuser.GoogleId = user.GoogleId
		dbuser.Name = user.Name
		dbuser.Email = user.Email

		if err := h.DB.Create(&dbuser).Error; err != nil {
			c.JSON(500, gin.H{"error" : "Internal server error"})
		}
	}

	myTokens, err := IssueToken(strconv.FormatUint(uint64(dbuser.ID), 10))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = Persist(context.Background(), h.Redis, myTokens)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	SendTokens(c, myTokens, dbuser)
}