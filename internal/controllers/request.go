package controllers

import (
	"math"
	"sahara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FetchRequests(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	category := c.Query("category")
	search := c.Query("search")
	location := c.Query("location")

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)

	var requests []models.Request

	offset := (p - 1) * l

	query := h.DB.Model(&models.Request{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if location != "" {
		query = query.Where("location = ?", location)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"(title LIKE ? OR description LIKE ?)",
			searchPattern,
			searchPattern,
		)
	}

	var total int64
	query.Count(&total)

	result := query.
		Limit(l).
		Offset(offset).
		Find(&requests)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{
		"page":       p,
		"limit":      l,
		"total":      total,
		"totalPages": math.Ceil(float64(total) / float64(l)),
		"data":       requests,
	})
}