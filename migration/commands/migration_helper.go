package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nahidhasan98/remind-name/config"
	"github.com/nahidhasan98/remind-name/name"
	"github.com/nahidhasan98/remind-name/platform"
	"go.mongodb.org/mongo-driver/bson"
)

// decodeAndConvert is a helper function that decodes JSON into a slice and converts it to []interface{}
func decodeAndConvert[T any](decoder *json.Decoder) ([]interface{}, error) {
	var items []T
	if err := decoder.Decode(&items); err != nil {
		return nil, err
	}

	documents := make([]interface{}, len(items))
	for i, item := range items {
		documents[i] = item
	}
	return documents, nil
}

// MigrateJSONToMongo is a generic helper function for JSON to MongoDB migration
func MigrateJSONToMongo(collectionName string, jsonFileName string) error {
	// Check if data already exists
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error checking %s: %v", collectionName, err)
	}

	if count > 0 {
		fmt.Printf("%s data already exists in database\n", collectionName)
		return nil
	}

	// Load data from JSON file
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", jsonFileName, err)
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	var documents []interface{}

	switch collectionName {
	case "platform":
		documents, err = decodeAndConvert[platform.Platform](decoder)
	case "name":
		documents, err = decodeAndConvert[name.Name](decoder)
	default:
		return fmt.Errorf("unsupported collection name: %s", collectionName)
	}

	if err != nil {
		return fmt.Errorf("error decoding %s: %v", jsonFileName, err)
	}

	// Insert data into MongoDB
	_, err = collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("error inserting %s: %v", collectionName, err)
	}

	fmt.Printf("Successfully migrated %s data to database\n", collectionName)
	return nil
}
