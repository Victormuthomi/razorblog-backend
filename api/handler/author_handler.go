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
// @Summary Register a new author
// @Description Creates a new author with hashed password
// @Tags Authors
// @Accept json
// @Produce json
// @Param author body object{ name=string, email=string, phone=string, password=string } true "Author info"
// @Success 201 {object} map[string]interface{} "Created author"
// @Failure 400 {object} map[string]string "bad request"
// @Failure 409 {object} map[string]string "email already registered"
// @Failure 500 {object} map[string]string "internal error"
// @Router /authors/register [post]
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

	author := &author.Author{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hashedPwd),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := h.Repo.CreateAuthor(author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Hide password in response
	created.Password = ""

	c.JSON(http.StatusCreated, map[string]interface{}{"author": created})
}

// LoginAuthor handles POST /authors/login
// @Summary Author login
// @Description Logs in an author and returns JWT token
// @Tags Authors
// @Accept json
// @Produce json
// @Param credentials body object{ email=string, password=string } true "Login info"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "bad request"
// @Failure 401 {object} map[string]string "invalid credentials"
// @Failure 500 {object} map[string]string "internal error"
// @Router /authors/login [post]
func (h *AuthorHandler) LoginAuthor(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	author, err := h.Repo.GetAuthorByEmail(req.Email)
	if err != nil || author == nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"author_id": author.ID.Hex(),
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
// @Summary Get author details
// @Description Retrieves author by ID
// @Tags Authors
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} map[string]interface{} "Author data"
// @Failure 400 {object} map[string]string "invalid id"
// @Failure 404 {object} map[string]string "author not found"
// @Security Bearer
// @Router /authors/{id} [get]
func (h *AuthorHandler) GetAuthor(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	author, err := h.Repo.GetAuthorByID(objID)
	if err != nil || author == nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": "author not found"})
		return
	}

	// Hide password in response
	author.Password = ""

	c.JSON(http.StatusOK, map[string]interface{}{"author": author})
}

// UpdateAuthor handles PUT /authors/:id
// @Summary Update author
// @Description Updates author data (password is hashed if provided)
// @Tags Authors
// @Accept json
// @Produce json
// @Param id path string true "Author ID"
// @Param author body object true "Updated author info"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "invalid id"
// @Failure 400 {object} map[string]string "bad request"
// @Failure 500 {object} map[string]string "internal error"
// @Security Bearer
// @Router /authors/{id} [put]
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
// @Summary Delete author
// @Description Deletes author by ID
// @Tags Authors
// @Produce json
// @Param id path string true "Author ID"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "invalid id"
// @Failure 500 {object} map[string]string "internal error"
// @Security Bearer
// @Router /authors/{id} [delete]
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

