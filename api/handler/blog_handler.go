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

// CreateBlog godoc
// @Summary Create a blog
// @Description Create a new blog post (protected)
// @Tags Blogs
// @Accept json
// @Produce json
// @Param blog body blog.Blog true "Blog payload"
// @Success 201 {object} blog.Blog
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /blogs [post]
// @Security Bearer
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

// GetBlog godoc
// @Summary Get a blog by ID
// @Description Get blog details by ID (public)
// @Tags Blogs
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} blog.Blog
// @Failure 404 {object} gin.H
// @Router /blogs/{id} [get]
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

// UpdateBlog godoc
// @Summary Update a blog
// @Description Update a blog post (protected, author only)
// @Tags Blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param blog body blog.Blog true "Blog payload"
// @Success 200 {object} blog.Blog
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /blogs/{id} [put]
// @Security Bearer
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

// DeleteBlog godoc
// @Summary Delete a blog
// @Description Delete a blog post (protected, author only)
// @Tags Blogs
// @Produce json
// @Param id path string true "Blog ID"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /blogs/{id} [delete]
// @Security Bearer
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

// ListBlogs godoc
// @Summary List blogs
// @Description List blogs with optional pagination (public)
// @Tags Blogs
// @Produce json
// @Param limit query int false "Limit"
// @Param skip query int false "Skip"
// @Success 200 {array} blog.Blog
// @Router /blogs [get]
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

