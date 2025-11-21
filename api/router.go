package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"razorblog-backend/api/handler"
	"razorblog-backend/api/middleware"
	"razorblog-backend/internal/repository"
)

// RegisterRoutes sets up all routes for the backend
func RegisterRoutes(r *gin.Engine, client *mongo.Client) {
	log.Println("Registering routes: / , /health, /authors, /blogs")

	// ===== Root Endpoint =====
	// @Summary Root endpoint
	// @Description Returns a welcome message for RazorBlog backend
	// @Tags Root
	// @Success 200 {object} map[string]string
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "RazorBlog Backend Running"})
	})
// HealthCheck godoc
// @Summary Health check
// @Description Checks if the server and MongoDB are running with retry
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health [get]
r.GET("/health", func(c *gin.Context) {
	const maxRetries = 3
	const retryDelay = time.Second

	dbStatus := "ok"
	httpStatus := http.StatusOK

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := client.Ping(ctx, nil)
		cancel()

		if err == nil {
			dbStatus = "ok"
			httpStatus = http.StatusOK
			break
		} else {
			dbStatus = "disconnected"
			httpStatus = http.StatusServiceUnavailable
			time.Sleep(retryDelay)
		}
	}

	c.JSON(httpStatus, gin.H{
		"server": "ok",
		"db":     dbStatus,
	})
})

		// ===== Author Routes =====
	db := client.Database("razorblog")
	authorRepo := repository.NewAuthorRepository(db)
	authorHandler := handler.NewAuthorHandler(authorRepo)

	// Public Author routes
	r.POST("/authors/register", authorHandler.RegisterAuthor)
	r.POST("/authors/login", authorHandler.LoginAuthor)
  r.GET("/authors/public/:id", authorHandler.GetPublicAuthor)


	// Protected Author routes
	authMiddleware := middleware.AuthMiddleware()
	authorProtected := r.Group("/authors", authMiddleware)
	{
		authorProtected.GET("/:id", authorHandler.GetAuthor)
		authorProtected.PUT("/:id", authorHandler.UpdateAuthor)
		authorProtected.DELETE("/:id", authorHandler.DeleteAuthor)
	}

	// ===== Blog Routes =====
	blogRepo := repository.NewBlogRepository(db)
  blogHandler := handler.NewBlogHandler(blogRepo, authorRepo) // pass authorRepo too


	// Public Blog routes
	r.GET("/blogs", blogHandler.ListBlogs)
  r.GET("/blogs/author/:author_id", blogHandler.GetBlogsByAuthor)
	r.GET("/blogs/:id", blogHandler.GetBlog)

	// Protected Blog routes
	blogProtected := r.Group("/blogs", authMiddleware)
	{
		// Create blog
		blogProtected.POST("", blogHandler.CreateBlog) // <- no trailing slash

		// Update blog
		blogProtected.PUT("/:id", blogHandler.UpdateBlog)

		// Delete blog
		blogProtected.DELETE("/:id", blogHandler.DeleteBlog)
	}

  
    // ❤️ New: Like/Unlike blog endpoints
    blogProtected.PATCH("/:id/like", blogHandler.LikeBlog)
    blogProtected.PATCH("/:id/unlike", blogHandler.UnlikeBlog)

	//  Comment routes
  // ===== Comment Routes =====
commentRepo := repository.NewCommentRepository(db)
commentHandler := handler.NewCommentHandler(commentRepo)

// Public Comment routes
// @Summary Create a comment
// @Description Create a new comment for a blog
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment payload"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "bad request"
// @Router /comments [post]
r.POST("/comments", commentHandler.CreateComment)

// @Summary List comments for a blog
// @Description List comments with optional pagination
// @Tags Comments
// @Produce json
// @Param blog_id path string true "Blog ID"
// @Param limit query int false "Limit"
// @Param skip query int false "Skip"
// @Success 200 {array} models.Comment
// @Failure 400 {object} map[string]string "invalid blog id"
// @Router /comments/{blog_id} [get]
r.GET("/comments/:blog_id", commentHandler.ListComments)

// @Summary Like a comment
// @Description Like a comment by username (one like per username)
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param body body object{username=string} true "Username liking the comment"
// @Success 200 {object} models.Comment
// @Failure 400 {object} map[string]string "bad request or already liked"
// @Router /comments/{id}/like [post]
r.POST("/comments/:id/like", commentHandler.LikeComment)

	// Share routes
  // ===== Share Routes =====
shareRepo := repository.NewShareRepository(db)
shareHandler := handler.NewShareHandler(shareRepo)

// Public Share routes
// @Summary Create a blog share
// @Description Record a share event for a blog
// @Tags Shares
// @Accept json
// @Produce json
// @Param share body models.Share true "Share payload"
// @Success 201 {object} models.Share
// @Failure 400 {object} map[string]string "bad request"
// @Router /shares [post]
r.POST("/shares", shareHandler.CreateShare)

// @Summary List shares for a blog
// @Description List all share events for a specific blog
// @Tags Shares
// @Produce json
// @Param blog_id path string true "Blog ID"
// @Success 200 {array} models.Share
// @Failure 400 {object} map[string]string "invalid blog id"
// @Router /shares/{blog_id} [get]
r.GET("/shares/:blog_id", shareHandler.ListShares)

}

