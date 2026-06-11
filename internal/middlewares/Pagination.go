package middlewares

import "github.com/gin-gonic/gin"


func DefaultPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		if query.Get("page") == "" {
			query.Set("page", "1")
		}

		if query.Get("limit") == "" {
			query.Set("limit", "10")
		}

		category := query.Get("category")

		if !(category == "education" || category == "healthcare" || category == "community" || category == "environment" || category == "disaster-relief" || category == "animal-welfare" || category == "others" || category == "") {
			query.Set("category", "")
		}

		c.Request.URL.RawQuery = query.Encode()

		c.Next()
	}
}