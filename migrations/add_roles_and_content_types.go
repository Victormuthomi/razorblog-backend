package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("razorblog")
	authors := db.Collection("authors")
	blogs := db.Collection("blogs")

	// 1. Migrate Authors to 'guest' role if missing
	authorFilter := bson.M{"role": bson.M{"$exists": false}}
	authorUpdate := bson.M{"$set": bson.M{"role": "guest"}}
	aResult, err := authors.UpdateMany(ctx, authorFilter, authorUpdate)
	if err != nil {
		log.Fatalf("Author migration failed: %v", err)
	}
	fmt.Printf("Authors: Matched %d, Modified %d\n", aResult.MatchedCount, aResult.ModifiedCount)

	// 2. Set Founder Role (Explicitly target your account)
	// Suggestion: Use your unique email from your .env or hardcode for this one-time run
	founderEmail := "your_email@example.com" 
	fResult, err := authors.UpdateOne(ctx, bson.M{"email": founderEmail}, bson.M{"$set": bson.M{"role": "founder"}})
	if err != nil {
		log.Printf("Founder update failed: %v", err)
	} else if fResult.MatchedCount > 0 {
		fmt.Println("Successfully assigned 'founder' role to:", founderEmail)
	}

	// 3. Migrate Blogs to 'blog' type if missing
	blogFilter := bson.M{"type": bson.M{"$exists": false}}
	blogUpdate := bson.M{"$set": bson.M{"type": "blog"}}
	bResult, err := blogs.UpdateMany(ctx, blogFilter, blogUpdate)
	if err != nil {
		log.Fatalf("Blog migration failed: %v", err)
	}
	fmt.Printf("Blogs: Matched %d, Modified %d\n", bResult.MatchedCount, bResult.ModifiedCount)
}
