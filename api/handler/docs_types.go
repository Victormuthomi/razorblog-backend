package handler

import (
	"razorblog-backend/internal/models/author"
)

// --- Author Requests ---

type AuthorRegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Phone    string `json:"phone" example:"+1234567890"`
	Password string `json:"password" example:"secret123"`
}

type AuthorLoginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secret123"`
}

// --- Author Responses ---

type AuthorResponse struct {
	Author author.Author `json:"author"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type MessageResponse struct {
	Message string `json:"message" example:"author updated"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

