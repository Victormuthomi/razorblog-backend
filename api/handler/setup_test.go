package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"razorblog-backend/internal/models/author"
	"razorblog-backend/internal/models/blog"
)

// --- MOCK BLOG REPO ---
type MockBlogRepo struct{ mock.Mock }

func (m *MockBlogRepo) Create(ctx context.Context, b *blog.Blog) (*blog.Blog, error) {
	args := m.Called(ctx, b)
	return args.Get(0).(*blog.Blog), args.Error(1)
}
func (m *MockBlogRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*blog.Blog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*blog.Blog), args.Error(1)
}
func (m *MockBlogRepo) Update(ctx context.Context, id primitive.ObjectID, u bson.M) (*blog.Blog, error) {
	args := m.Called(ctx, id, u)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*blog.Blog), args.Error(1)
}
func (m *MockBlogRepo) Delete(ctx context.Context, id primitive.ObjectID) error { return m.Called(ctx, id).Error(0) }
func (m *MockBlogRepo) List(ctx context.Context, l, s int64) ([]*blog.Blog, error) { return nil, nil }
func (m *MockBlogRepo) ListByAuthor(ctx context.Context, id primitive.ObjectID) ([]*blog.Blog, error) { return nil, nil }
func (m *MockBlogRepo) IncrementReaders(ctx context.Context, id primitive.ObjectID) error { return nil }
func (m *MockBlogRepo) LikeBlog(ctx context.Context, bID, uID primitive.ObjectID) error { return nil }
func (m *MockBlogRepo) UnlikeBlog(ctx context.Context, bID, uID primitive.ObjectID) error { return nil }

// --- MOCK AUTHOR REPO ---
type MockAuthorRepo struct{ mock.Mock }

func (m *MockAuthorRepo) GetAuthorByID(id primitive.ObjectID) (*author.Author, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*author.Author), args.Error(1)
}
func (m *MockAuthorRepo) GetAuthorByEmail(e string) (*author.Author, error) {
	args := m.Called(e)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*author.Author), args.Error(1)
}
func (m *MockAuthorRepo) CreateAuthor(a *author.Author) (*author.Author, error) {
	args := m.Called(a)
	return args.Get(0).(*author.Author), args.Error(1)
}
func (m *MockAuthorRepo) UpdateAuthor(id primitive.ObjectID, u bson.M) error {
	return m.Called(id, u).Error(0)
}
func (m *MockAuthorRepo) DeleteAuthor(id primitive.ObjectID) error { return m.Called(id).Error(0) }
