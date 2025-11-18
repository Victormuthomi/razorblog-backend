package database

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.TODO()

func Connect(uri string) (*mongo.Client, error) {
    clientOptions := options.Client().ApplyURI(uri)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    // Ping to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        return nil, err
    }

    return client, nil
}

