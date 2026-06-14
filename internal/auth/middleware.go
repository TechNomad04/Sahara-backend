package auth

import (
	"context"
	"net/http"
	"sahara/internal/store"
	"strings"

	"github.com/gin-gonic/gin"
)

func BearerFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func AuthMiddleware(r *store.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenStr := BearerFromHeader(c)

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing bearer token",
			})
			return
		}

		claims, err := ParseAccess(tokenStr)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		ctx := context.Background()

		if _, err := r.GetUserByJTI(
			ctx,
			"access:"+claims.ID,
		); err != nil {

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token revoked",
			})
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("entityType", claims.EntityType)

		c.Next()
	}
}