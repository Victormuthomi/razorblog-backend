package blog

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Blog struct ...
type Blog struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	AuthorID  primitive.ObjectID   `bson:"author_id" json:"author_id"`
	Title     string               `bson:"title" json:"title"`
	Content   string               `bson:"content" json:"content"`
	ImageURL  string               `bson:"image_url,omitempty" json:"image_url"`
	Category  string               `bson:"category" json:"category"`
	Readers   int                  `bson:"readers" json:"readers"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`

	// ❤️ New feature: store user IDs who liked this blog
	Likes []primitive.ObjectID `bson:"likes,omitempty" json:"likes,omitempty"`
}

