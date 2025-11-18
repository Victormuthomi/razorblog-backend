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

// CreateComment godoc
// @Summary Create a comment
// @Description Create a new comment for a blog
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment payload"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "bad request"
// @Router /comments [post]
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

// ListComments godoc
// @Summary List comments for a blog
// @Description List comments for a given blog with optional pagination
// @Tags Comments
// @Produce json
// @Param blog_id path string true "Blog ID"
// @Param limit query int false "Limit (default 10)"
// @Param skip query int false "Skip (default 0)"
// @Success 200 {array} models.Comment
// @Failure 400 {object} map[string]string "invalid blog id"
// @Router /comments/{blog_id} [get]
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

// LikeComment godoc
// @Summary Like a comment
// @Description Like a comment by username. One like per username only.
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param body body object{username=string} true "Username liking the comment"
// @Success 200 {object} models.Comment
// @Failure 400 {object} map[string]string "bad request or already liked"
// @Router /comments/{id}/like [post]
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

