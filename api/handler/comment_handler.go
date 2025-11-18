package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"razorblog-backend/internal/repository"
	models "razorblog-backend/internal/models/comment"
)

// CommentHandler handles HTTP requests for comments
type CommentHandler struct {
	repo *repository.CommentRepository
}

func NewCommentHandler(repo *repository.CommentRepository) *CommentHandler {
	return &CommentHandler{repo: repo}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var cmt models.Comment
	if err := c.ShouldBindJSON(&cmt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cmt.BlogID.IsZero() || cmt.Username == "" || cmt.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "blog_id, username, and content are required"})
		return
	}

	created, err := h.repo.Create(context.Background(), &cmt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *CommentHandler) ListComments(c *gin.Context) {
	blogIDStr := c.Param("blog_id")
	blogID, err := primitive.ObjectIDFromHex(blogIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog ID"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)

	comments, err := h.repo.List(context.Background(), blogID, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	var body struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	updated, err := h.repo.Like(context.Background(), commentID, body.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

