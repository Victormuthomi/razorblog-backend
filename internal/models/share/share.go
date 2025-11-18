package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// Share represents a blog share event (optional tracking)
type Share struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    BlogID    primitive.ObjectID `bson:"blog_id" json:"blog_id"`       // Blog being shared
    Platform  string             `bson:"platform" json:"platform"`     // e.g., Twitter, Facebook
    CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Timestamp of the share
}

