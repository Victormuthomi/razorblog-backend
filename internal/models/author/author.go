package author

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// Author represents a blogger who can write blogs
type Author struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`          // MongoDB ObjectID
    Name      string             `bson:"name" json:"name"`                 // Full name of the author
    Email     string             `bson:"email" json:"email"`               // Email for login/auth
    Password  string             `bson:"password" json:"-"`                // Hashed password, excluded from JSON responses
    Phone     string             `bson:"phone,omitempty" json:"phone"`     // Optional phone number
    AvatarURL string             `bson:"avatar_url,omitempty" json:"avatar_url"` // URL to author's avatar image
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`     // Timestamp when author was created
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`     // Timestamp when author was last updated
}

