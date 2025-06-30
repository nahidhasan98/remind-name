package name

import (
	"github.com/nahidhasan98/remind-name/config"
	"github.com/nahidhasan98/remind-name/logger"
	"go.mongodb.org/mongo-driver/bson"
)

type repository struct {
	Collection string
}

func newRepository() *repository {
	return &repository{
		Collection: "name",
	}
}

func (repo *repository) getName(id int) (*Name, error) {
	// connecting to DB
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	var name Name
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&name)
	if err != nil {
		logger.Error("Failed to get name for id %d: %v", id, err)
		return nil, err
	}

	logger.Debug("Fetched name for id %d: %+v", id, name)
	return &name, nil
}
