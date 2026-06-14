package auth

import (
	"context"
	"encoding/json"
	"errors"
	"sahara/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	token, err := GoogleConfig.Exchange(
		context.Background(),
		code,
	)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
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
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer resp.Body.Close()

	var googleUser struct {
		GoogleId string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch state {

	case "user":

		var dbuser models.User

		err := h.DB.
			Where("google_id = ?", googleUser.GoogleId).
			First(&dbuser).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				dbuser = models.User{
					GoogleId: googleUser.GoogleId,
					Name:     googleUser.Name,
					Email:    googleUser.Email,
				}

				if err := h.DB.Create(&dbuser).Error; err != nil {
					c.JSON(500, gin.H{
						"error": "internal server error",
					})
					return
				}

			} else {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		tokens, err := IssueToken(
			strconv.FormatUint(uint64(dbuser.ID), 10),
			"user",
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := Persist(
			context.Background(),
			h.Redis,
			tokens,
		); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		SendTokens(c, tokens, dbuser)

	case "organization":

		var org models.Organization

		err := h.DB.
			Where("google_id = ?", googleUser.GoogleId).
			First(&org).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				org = models.Organization{
					GoogleId: googleUser.GoogleId,
					Name:     googleUser.Name,
					Email:    googleUser.Email,
				}

				if err := h.DB.Create(&org).Error; err != nil {
					c.JSON(500, gin.H{
						"error": "internal server error",
					})
					return
				}

			} else {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
		}

		tokens, err := IssueToken(
			strconv.FormatUint(uint64(org.ID), 10),
			"organization",
		)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := Persist(
			context.Background(),
			h.Redis,
			tokens,
		); err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		SendTokens(c, tokens, org)

	default:
		c.JSON(400, gin.H{
			"error": "invalid state",
		})
	}
}