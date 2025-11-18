package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"razorblog-backend/internal/repository"
	models "razorblog-backend/internal/models/share"
)

// ShareHandler handles HTTP requests for blog shares
type ShareHandler struct {
	repo *repository.ShareRepository
}

func NewShareHandler(repo *repository.ShareRepository) *ShareHandler {
	return &ShareHandler{repo: repo}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	var s models.Share
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if s.BlogID.IsZero() || s.Platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "blog_id and platform are required"})
		return
	}

	created, err := h.repo.Create(context.Background(), &s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create share"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ShareHandler) ListShares(c *gin.Context) {
	blogIDStr := c.Param("blog_id")
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}

	shares, err := h.repo.List(context.Background(), blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shares"})
		return
	}

	c.JSON(http.StatusOK, shares)
}

