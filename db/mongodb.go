package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var dbUrl = "mongodb://localhost:27017"

func ConnectDB() *mongo.Collection {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!", client, ctx)
	collection := client.Database("demo").Collection("books")
	fmt.Println("Collection: ", collection)
	return collection
}
