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
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("razorblog")
	authors := db.Collection("authors")

	// Update all authors missing 'bio'
	update := bson.M{"$set": bson.M{"bio": ""}}
	filter := bson.M{"bio": bson.M{"$exists": false}}

	result, err := authors.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %d authors\n", result.MatchedCount)
	fmt.Printf("Modified %d authors\n", result.ModifiedCount)
}

