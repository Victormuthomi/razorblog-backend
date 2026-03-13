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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"razorblog-backend/internal/models/author"
)

func TestAuthor_Security(t *testing.T) {
	gin.SetMode(gin.TestMode)

t.Run("FORCE GUEST: Registration ignores role in JSON", func(t *testing.T) {
    mAuth := new(MockAuthorRepo)
    h := NewAuthorHandler(mAuth)

    var capturedAuthor *author.Author
    // Ensure we return an empty author struct on success so pointers aren't nil
    mAuth.On("GetAuthorByEmail", mock.Anything).Return((*author.Author)(nil), nil)
    mAuth.On("CreateAuthor", mock.Anything).Run(func(args mock.Arguments) {
        capturedAuthor = args.Get(0).(*author.Author)
    }).Return(&author.Author{}, nil)

    w := httptest.NewRecorder()
    _, r := gin.CreateTestContext(w)
    r.POST("/register", h.RegisterAuthor)

    // USE A STRUCT: This ensures JSON types (like strings vs numbers) are correct
    payload := struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }{
        Name:     "Hacker",
        Email:    "hacker@test.com",
        Password: "securepassword123",
        Role:     "founder",
    }
    
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
    r.ServeHTTP(w, req)

    // CHECK STATUS BEFORE CHECKING POINTERS
    if !assert.Equal(t, http.StatusCreated, w.Code, "Registration failed with: "+w.Body.String()) {
        t.FailNow() // Stop here if status is 400 so we don't panic below
    }
    
    assert.NotNil(t, capturedAuthor)
    assert.Equal(t, author.RoleGuest, capturedAuthor.Role)
})	
	t.Run("PROTECT UPDATE: User cannot inject role field", func(t *testing.T) {
		mAuth := new(MockAuthorRepo)
		h := NewAuthorHandler(mAuth)

		userID := primitive.NewObjectID()
		var capturedUpdate bson.M
		mAuth.On("UpdateAuthor", userID, mock.Anything).Run(func(args mock.Arguments) {
			capturedUpdate = args.Get(1).(bson.M)
		}).Return(nil)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.PUT("/authors/:id", h.UpdateAuthor)

		updatePayload := map[string]interface{}{"name": "New Name", "role": "founder"}
		body, _ := json.Marshal(updatePayload)
		req, _ := http.NewRequest("PUT", "/authors/"+userID.Hex(), bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		_, roleExists := capturedUpdate["role"]
		assert.False(t, roleExists)
	})
}
