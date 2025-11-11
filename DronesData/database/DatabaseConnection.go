package database

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	clientOnce sync.Once
)

// GetClient returns a singleton MongoDB client instance
func GetClient() *mongo.Client {
	clientOnce.Do(func() {
		uri := os.Getenv("DB_URI")
		if uri == "" {
			log.Fatal("DB URL coudln't load from properties.env file")
		}

		clientOptions := options.Client().
			ApplyURI(uri).
			SetConnectTimeout(10 * time.Second).
			SetServerSelectionTimeout(5 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}

		// Verify connection
		err = c.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to ping DB: %v", err)
		}

		client = c
		log.Println("DB connection established")
	})
	return client
}
