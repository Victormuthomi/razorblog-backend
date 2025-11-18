package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Comment represents a comment left by a reader on a blog
type Comment struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	BlogID    primitive.ObjectID   `bson:"blog_id" json:"blog_id"`             // Blog this comment belongs to
	Username  string               `bson:"username" json:"username"`           // Name of the commentor
	Content   string               `bson:"content" json:"content"`             // Comment text
	Likes     int                  `bson:"likes" json:"likes"`                 // Number of likes
	LikedBy   []string             `bson:"liked_by,omitempty" json:"liked_by"` // Track users who liked this comment (to ensure one like per person)
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`       // Timestamp
}

