package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"razorblog-backend/internal/models/author"
	"razorblog-backend/internal/models/blog"
)

type IBlogRepository interface {
	Create(ctx context.Context, b *blog.Blog) (*blog.Blog, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*blog.Blog, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*blog.Blog, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit int64, skip int64) ([]*blog.Blog, error)
	ListByAuthor(ctx context.Context, authorID primitive.ObjectID) ([]*blog.Blog, error)
	IncrementReaders(ctx context.Context, id primitive.ObjectID) error
	LikeBlog(ctx context.Context, blogID, userID primitive.ObjectID) error
	UnlikeBlog(ctx context.Context, blogID, userID primitive.ObjectID) error
}

type IAuthorRepository interface {
	CreateAuthor(a *author.Author) (*author.Author, error)
	GetAuthorByID(id primitive.ObjectID) (*author.Author, error)
	GetAuthorByEmail(email string) (*author.Author, error)
	UpdateAuthor(id primitive.ObjectID, update bson.M) error
	DeleteAuthor(id primitive.ObjectID) error
}
