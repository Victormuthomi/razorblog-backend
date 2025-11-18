package handler

import (
    "net/http"
    "time"

  "razorblog-backend/internal/models/author"
	"razorblog-backend/internal/repository"


    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson/primitive"
)


// JWT secret (load from env in production)
var jwtSecret = []byte("supersecretkey")

// AuthorHandler holds repository reference
type AuthorHandler struct {
	Repo *repository.AuthorRepository
}

// NewAuthorHandler creates a new AuthorHandler
func NewAuthorHandler(repo *repository.AuthorRepository) *AuthorHandler {
	return &AuthorHandler{Repo: repo}
}

// RegisterAuthor handles POST /authors/register
func (h *AuthorHandler) RegisterAuthor(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Check if author already exists
	existing, _ := h.Repo.GetAuthorByEmail(req.Email)
	if existing != nil {
		c.JSON(http.StatusConflict, map[string]string{"error": "email already registered"})
		return
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
		return
	}

	newAuthor := &author.Author{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := h.Repo.CreateAuthor(newAuthor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Hide password in response
	created.Password = ""

	c.JSON(http.StatusCreated, map[string]interface{}{"author": created})
}

// LoginAuthor handles POST /authors/login
func (h *AuthorHandler) LoginAuthor(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	authorObj, err := h.Repo.GetAuthorByEmail(req.Email)
	if err != nil || authorObj == nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(authorObj.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"author_id": authorObj.ID.Hex(),
		"exp":       time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

// GetAuthor handles GET /authors/:id
func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	authorObj, err := h.Repo.GetAuthorByID(objID)
	if err != nil || authorObj == nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "author not found"})
		return
	}

	// Hide password in response
	authorObj.Password = ""

	c.JSON(http.StatusOK, map[string]interface{}{"author": authorObj})
}

// UpdateAuthor handles PUT /authors/:id
func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// If password is being updated, hash it
	if pwd, ok := update["password"].(string); ok && pwd != "" {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
			return
		}
		update["password"] = string(hashedPwd)
	}

	update["updated_at"] = time.Now()

	if err := h.Repo.UpdateAuthor(objID, update); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": "author updated"})
}

// DeleteAuthor handles DELETE /authors/:id
func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := h.Repo.DeleteAuthor(objID); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": "author deleted"})
}

