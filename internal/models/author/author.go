package author

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type UserRole string

const (
    RoleFounder UserRole = "founder"
    RoleGuest   UserRole = "guest"
)

type Author struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Password  string             `bson:"password" json:"-"`
    Role      UserRole           `bson:"role" json:"role"` // New: founder vs guest
    Phone     string             `bson:"phone,omitempty" json:"phone"`
    AvatarURL string             `bson:"avatar_url,omitempty" json:"avatar_url"`
    Bio       string             `bson:"bio,omitempty" json:"bio"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
