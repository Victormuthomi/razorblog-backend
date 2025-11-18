package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"razorblog-backend/internal/models/blog"
	"razorblog-backend/internal/repository"
)

type BlogHandler struct {
	repo *repository.BlogRepository
}

func NewBlogHandler(repo *repository.BlogRepository) *BlogHandler {
	return &BlogHandler{repo: repo}
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var b blog.Blog
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, exists := c.Get("author_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(authorID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid author ID"})
		return
	}
	b.AuthorID = objID

	created, err := h.repo.Create(context.Background(), &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *BlogHandler) GetBlog(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}

	// Increment readers count automatically
	_ = h.repo.IncrementReaders(context.Background(), objID)

	b, err := h.repo.GetByID(context.Background(), objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	c.JSON(http.StatusOK, b)
}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}

	var b blog.Blog
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := map[string]interface{}{
		"title":     b.Title,
		"content":   b.Content,
		"image_url": b.ImageURL,
		"category":  b.Category,
	}

	updated, err := h.repo.Update(context.Background(), objID, update)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}

	if err := h.repo.Delete(context.Background(), objID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog deleted"})
}

func (h *BlogHandler) ListBlogs(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	blogs, err := h.repo.List(context.Background(), limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blogs)
}

