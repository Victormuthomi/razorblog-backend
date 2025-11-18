package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"razorblog-backend/internal/models/blog"
)

type BlogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		collection: db.Collection("blogs"),
	}
}

func (r *BlogRepository) Create(ctx context.Context, b *blog.Blog) (*blog.Blog, error) {
	b.ID = primitive.NewObjectID()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	b.Readers = 0
	_, err := r.collection.InsertOne(ctx, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BlogRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*blog.Blog, error) {
	var b blog.Blog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BlogRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*blog.Blog, error) {
	update["updated_at"] = time.Now()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedBlog blog.Blog
	err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": update}, opts).Decode(&updatedBlog)
	if err != nil {
		return nil, err
	}
	return &updatedBlog, nil
}

func (r *BlogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *BlogRepository) List(ctx context.Context, limit int64, skip int64) ([]*blog.Blog, error) {
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []*blog.Blog
	for cursor.Next(ctx) {
		var b blog.Blog
		if err := cursor.Decode(&b); err != nil {
			return nil, err
		}
		blogs = append(blogs, &b)
	}
	return blogs, nil
}

func (r *BlogRepository) IncrementReaders(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"readers": 1}})
	return err
}

