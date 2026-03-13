package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"razorblog-backend/internal/models/author"
	"razorblog-backend/internal/models/blog"
)

func TestCreateBlog_Security(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("REJECT: Guest attempts to post TDD content", func(t *testing.T) {
		mBlog := new(MockBlogRepo)
		mAuth := new(MockAuthorRepo)
		h := NewBlogHandler(mBlog, mAuth)

		guestID := primitive.NewObjectID()
		mAuth.On("GetAuthorByID", guestID).Return(&author.Author{Role: "guest"}, nil)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		r.POST("/blogs", func(ctx *gin.Context) {
			ctx.Set("author_id", guestID.Hex())
			h.CreateBlog(ctx)
		})

		input := blog.Blog{Title: "Guest Specs", Type: blog.TypeTDD}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/blogs", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mBlog.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("ALLOW: Founder posts TDD content", func(t *testing.T) {
		mBlog := new(MockBlogRepo)
		mAuth := new(MockAuthorRepo)
		h := NewBlogHandler(mBlog, mAuth)

		founderID := primitive.NewObjectID()
		mAuth.On("GetAuthorByID", founderID).Return(&author.Author{Role: "founder"}, nil)
		mBlog.On("Create", mock.Anything, mock.Anything).Return(&blog.Blog{}, nil)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		r.POST("/blogs", func(ctx *gin.Context) {
			ctx.Set("author_id", founderID.Hex())
			h.CreateBlog(ctx)
		})

		input := blog.Blog{Title: "Founder Architecture", Type: blog.TypeTDD}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/blogs", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
