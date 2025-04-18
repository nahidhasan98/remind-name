package name

import (
	"github.com/nahidhasan98/remind-name/config"
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
		return nil, err
	}

	return &name, nil
}
