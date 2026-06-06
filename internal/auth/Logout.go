package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler)Logout (c *gin.Context) {
	type LogoutRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error" : "Invalid JSON",
		})

		return
	}

	access_token := BearerFromHeader(c)

	if access_token == "" {
		c.JSON(http.StatusUnauthorized, gin.H {
			"Error" : "missing access token",
		})
		return
	}

	ctx := context.Background()

	access_claims, err := ParseAccess(access_token)

	if err == nil {

		_ = h.Redis.DelJTI(
			ctx,
			"access:"+access_claims.ID,
		)
	}

	refreshClaims, err := ParseRefresh(req.RefreshToken)

	if err == nil {

		_ = h.Redis.DelJTI(
			ctx,
			"refresh:"+refreshClaims.ID,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}