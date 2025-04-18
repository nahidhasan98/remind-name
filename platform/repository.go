package platform

import (
	"github.com/nahidhasan98/remind-name/config"
	"go.mongodb.org/mongo-driver/bson"
)

type repository struct {
	Collection string
}

func newRepository() *repository {
	return &repository{
		Collection: "platform",
	}
}

func (repo *repository) getPlatformDetailsByName(name string) (*Platform, error) {
	// connecting to DB
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	var platform Platform
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&platform)
	if err != nil {
		return nil, err
	}

	return &platform, nil
}

func (repo *repository) getAllPlatforms() ([]Platform, error) {
	// connecting to DB
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	var platform []Platform
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &platform)
	if err != nil {
		return nil, err
	}

	return platform, nil
}
