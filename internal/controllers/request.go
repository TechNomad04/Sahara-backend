package controllers

import (
	"net/http"
	"sahara/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FetchRequests(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	category := c.Query("category")
	search := c.Query("search")
	country := c.Query("country")
	state := c.Query("state")
	city := c.Query("city")

	p, _ := strconv.Atoi(page)
	l, _ := strconv.Atoi(limit)

	type RequestResponse struct {
		ID               uint     `json:"id"`
		Title            string   `json:"title"`
		Description      string   `json:"description"`
		Categories       []string `json:"categories"`
		Location         string   `json:"location"`
		Status           string   `json:"status"`
		IssuerName       string   `json:"issuerName"`
		IssuerType       string   `json:"issuerType"`
		ParticipantCount int      `json:"participantCount"`
		CreatedAt        string   `json:"createdAt"`
	}

	offset := (p - 1) * l

	query := h.DB.Model(&models.Request{})

	if category != "" {
		query = query.Where("? = ANY(categories)", category)
	}

	if country != "" {
		query = query.Where("requests.country = ?", country)
	}

	if state != "" {
		query = query.Where("requests.state = ?", state)
	}

	if city != "" {
		query = query.Where("requests.city = ?", city)
	}

	if search != "" {
		searchPattern := "%" + search + "%"

		query = query.Where(
			"(requests.title ILIKE ? OR requests.description ILIKE ?)",
			searchPattern,
			searchPattern,
		)
	}

	var total int64
	query.Count(&total)

	var requests []models.Request

	err := query.
		Preload("Issuer").
		Limit(l).
		Offset(offset).
		Order("requests.created_at DESC").
		Find(&requests).Error

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]RequestResponse, 0, len(requests))

	for _, req := range requests {
		response = append(response, RequestResponse{
			ID:               req.ID,
			Title:            req.Title,
			Description:      req.Description,
			Categories:       []string(req.Categories),
			Location:         req.City + ", " + req.State + ", " + req.Country,
			Status:           string(req.Status),
			IssuerName:       req.Issuer.Name,
			IssuerType:       req.Issuer.IssuerType,
			ParticipantCount: req.ParticipantCount,
			CreatedAt:        req.CreatedAt.Format("Jan 02, 2006"),
		})
	}

	c.JSON(200, gin.H{
		"requests": response,
		"hasMore":  offset+l < int(total),
		"total":    total,
	})
}

func (h *Handler) CreateRequest(c *gin.Context) {

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
		Country     string   `json:"country"`
		State       string   `json:"state"`
		City        string   `json:"city"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	entityType := c.MustGet("entityType").(string)

	if entityType != "organization" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "only organizations can create requests",
		})
		return
	}

	userID := c.MustGet("userID").(string)

	orgID, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid organization id",
		})
		return
	}

	request := models.Request{
		Title:          req.Title,
		Description:    req.Description,
		Categories:     req.Categories,
		IssuerId: uint(orgID),
		Country:        req.Country,
		State:          req.State,
		City:           req.City,
	}

	// TODO:
	// semantic search duplicate detection

	if err := h.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create request",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"request": request,
	})
}
