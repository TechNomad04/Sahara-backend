package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(c *gin.Context) {

	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})

		return
	}

	if req.RefreshToken == "" {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing refresh token",
		})

		return
	}
	

	claims, err := ParseRefresh(
		req.RefreshToken,
	)

	if err != nil {

		fmt.Println("PARSE ERROR:", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx := context.Background()

	_, err = h.Redis.GetUserByJTI(
		ctx,
		"refresh:"+claims.ID,
	)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token revoked",
		})

		return
	}

	_ = h.Redis.DelJTI(
		ctx,
		"refresh:"+claims.ID,
	)

	toks, err := IssueToken(
		claims.Subject,
		claims.EntityType,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not issue new tokens",
		})

		return
	}

	err = Persist(
		ctx,
		h.Redis,
		toks,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not persist new tokens",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{

		"access_token":  toks.Access,
		"refresh_token": toks.Refresh,

		"access_exp":  toks.ExpAcc,
		"refresh_exp": toks.ExpRef,
	})
}