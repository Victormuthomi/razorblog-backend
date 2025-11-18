package repository

import (
    "context"
    "errors"
    "time"

    "razorblog-backend/internal/models/author"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// AuthorRepository manages CRUD operations for Author
type AuthorRepository struct {
    collection *mongo.Collection
}

// NewAuthorRepository returns a new AuthorRepository instance
func NewAuthorRepository(db *mongo.Database) *AuthorRepository {
    return &AuthorRepository{
        collection: db.Collection("authors"),
    }
}

// CreateAuthor inserts a new author into the DB
func (r *AuthorRepository) CreateAuthor(a *author.Author) (*author.Author, error) {
    now := time.Now()
    a.ID = primitive.NewObjectID()
    a.CreatedAt = now
    a.UpdatedAt = now

    _, err := r.collection.InsertOne(context.Background(), a)
    if err != nil {
        return nil, err
    }
    return a, nil
}

// GetAuthorByID finds an author by ID
func (r *AuthorRepository) GetAuthorByID(id primitive.ObjectID) (*author.Author, error) {
    var a author.Author
    err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&a)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }
        return nil, err
    }
    return &a, nil
}

// GetAuthorByEmail finds an author by email (useful for login)
func (r *AuthorRepository) GetAuthorByEmail(email string) (*author.Author, error) {
    var a author.Author
    err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&a)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return nil, nil
        }
        return nil, err
    }
    return &a, nil
}

// UpdateAuthor updates an existing author
func (r *AuthorRepository) UpdateAuthor(id primitive.ObjectID, update bson.M) error {
    update["updated_at"] = time.Now()
    _, err := r.collection.UpdateOne(
        context.Background(),
        bson.M{"_id": id},
        bson.M{"$set": update},
    )
    return err
}

// DeleteAuthor removes an author by ID
func (r *AuthorRepository) DeleteAuthor(id primitive.ObjectID) error {
    _, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
    return err
}

