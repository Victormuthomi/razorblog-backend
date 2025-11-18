package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	models "razorblog-backend/internal/models/share"
)

// ShareRepository handles database operations for blog shares
type ShareRepository struct {
	collection *mongo.Collection
}

func NewShareRepository(db *mongo.Database) *ShareRepository {
	return &ShareRepository{
		collection: db.Collection("shares"),
	}
}

// Create inserts a new share into the database
func (r *ShareRepository) Create(ctx context.Context, s *models.Share) (*models.Share, error) {
	s.ID = primitive.NewObjectID()
	s.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// List returns all shares for a given blog
func (r *ShareRepository) List(ctx context.Context, blogID primitive.ObjectID) ([]*models.Share, error) {
	cursor, err := r.collection.Find(ctx, map[string]interface{}{"blog_id": blogID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shares []*models.Share
	for cursor.Next(ctx) {
		var s models.Share
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		shares = append(shares, &s)
	}

	return shares, nil
}

