package author

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// Author represents a blogger who can write blogs
type Author struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"-"`       // hide email for public
    Password  string             `bson:"password" json:"-"`
    Phone     string             `bson:"phone,omitempty" json:"-"` // hide phone
    AvatarURL string             `bson:"avatar_url,omitempty" json:"avatar_url"`
    Bio       string             `bson:"bio,omitempty" json:"bio"` // new bio field
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

