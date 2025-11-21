package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"razorblog-backend/internal/models/blog"
	"razorblog-backend/internal/models/author"
)

type BlogRepository struct {
	collection    *mongo.Collection
	authorCol     *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		collection: db.Collection("blogs"),
		authorCol:  db.Collection("authors"),
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

// New method to get a blog along with author's name
func (r *BlogRepository) GetBlogWithAuthor(ctx context.Context, id primitive.ObjectID) (*blog.Blog, string, error) {
	b, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, "", err
	}

	var a author.Author
	if err := r.authorCol.FindOne(ctx, bson.M{"_id": b.AuthorID}).Decode(&a); err != nil {
		return b, "", nil // author might not exist, still return blog
	}

	return b, a.Name, nil
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

// New method to list blogs with author names
func (r *BlogRepository) GetBlogsWithAuthors(ctx context.Context, limit, skip int64) ([]*blog.Blog, map[primitive.ObjectID]string, error) {
	blogs, err := r.List(ctx, limit, skip)
	if err != nil {
		return nil, nil, err
	}

	authorIDs := make([]primitive.ObjectID, 0, len(blogs))
	for _, b := range blogs {
		authorIDs = append(authorIDs, b.AuthorID)
	}

	cursor, err := r.authorCol.Find(ctx, bson.M{"_id": bson.M{"$in": authorIDs}})
	if err != nil {
		return blogs, nil, nil
	}
	defer cursor.Close(ctx)

	authorMap := make(map[primitive.ObjectID]string)
	for cursor.Next(ctx) {
		var a author.Author
		if err := cursor.Decode(&a); err == nil {
			authorMap[a.ID] = a.Name
		}
	}

	return blogs, authorMap, nil
}

func (r *BlogRepository) IncrementReaders(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"readers": 1}})
	return err
}

func (r *BlogRepository) LikeBlog(ctx context.Context, blogID, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{"$addToSet": bson.M{"likes": userID}},
	)
	return err
}

func (r *BlogRepository) UnlikeBlog(ctx context.Context, blogID, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{"$pull": bson.M{"likes": userID}},
	)
	return err
}

//Getting a blog by author Id
func (r *BlogRepository) ListByAuthor(ctx context.Context, authorID primitive.ObjectID) ([]*blog.Blog, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"author_id": authorID},
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
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

