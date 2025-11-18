package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "razorblog-backend/internal/models/comment"
)

// CommentRepository handles database operations for comments
type CommentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		collection: db.Collection("comments"),
	}
}

// Create inserts a new comment into the database
func (r *CommentRepository) Create(ctx context.Context, cmt *models.Comment) (*models.Comment, error) {
	cmt.ID = primitive.NewObjectID()
	cmt.CreatedAt = time.Now()
	cmt.Likes = 0
	cmt.LikedBy = []string{}
	_, err := r.collection.InsertOne(ctx, cmt)
	if err != nil {
		return nil, err
	}
	return cmt, nil
}

// List returns comments for a specific blog with pagination
func (r *CommentRepository) List(ctx context.Context, blogID primitive.ObjectID, limit, skip int64) ([]*models.Comment, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(ctx, bson.M{"blog_id": blogID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*models.Comment
	for cursor.Next(ctx) {
		var c models.Comment
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}

// Like adds a like to a comment if the username hasn't liked it yet
func (r *CommentRepository) Like(ctx context.Context, commentID primitive.ObjectID, username string) (*models.Comment, error) {
	// Check if user already liked
	var c models.Comment
	if err := r.collection.FindOne(ctx, bson.M{"_id": commentID}).Decode(&c); err != nil {
		return nil, err
	}

	for _, user := range c.LikedBy {
		if user == username {
			return nil, errors.New("user already liked this comment")
		}
	}

	update := bson.M{
		"$inc":  bson.M{"likes": 1},
		"$push": bson.M{"liked_by": username},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedComment models.Comment
	if err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": commentID}, update, opts).Decode(&updatedComment); err != nil {
		return nil, err
	}

	return &updatedComment, nil
}

// Delete removes a comment by ID
func (r *CommentRepository) Delete(ctx context.Context, commentID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": commentID})
	return err
}

