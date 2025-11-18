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

// CreateShare godoc
// @Summary Create a new share
// @Description Records a share of a blog post on a platform
// @Tags Shares
// @Accept json
// @Produce json
// @Param share body map[string]string true "Share info (blog_id, platform)"
// @Success 201 {object} models.Share
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shares [post]
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

// ListShares godoc
// @Summary List shares for a blog
// @Description Returns a list of shares for a specific blog post
// @Tags Shares
// @Produce json
// @Param blog_id path string true "Blog ID"
// @Success 200 {array} models.Share
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /blogs/{blog_id}/shares [get]
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

