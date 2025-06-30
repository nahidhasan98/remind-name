package config

import (
	"context"
	"time"

	"github.com/nahidhasan98/remind-name/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DBConnect() (*mongo.Database, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_CONNECTION_STRING))
	if err != nil {
		logger.Error("database initialization error: %v", err)
	}

	if err = dbClient.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error("database initialization error: %v", err)
	}

	return dbClient.Database(DB_NAME), ctx, cancel
}

// Create necessary indexes
func createIndexes(db *mongo.Database, collectionName string) error {
	collection := db.Collection(collectionName)

	// Define the index models
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.M{"time_from": 1}, // Index for optimizing time-based queries
		},
		{
			Keys: bson.M{"time_to": 1}, // Index for filtering users based on start time
		},
		{
			Keys: bson.M{"last_sent_at": 1}, // Index for filtering users based on end time
		},
	}

	// Create indexes
	_, err := collection.Indexes().CreateMany(context.TODO(), indexModels)
	return err
}

func init() {
	DB, ctx, cancel := DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	if err := createIndexes(DB, SUB_COLLECTION); err != nil {
		logger.Error("error creating indexes: %v", err)
	}
}
