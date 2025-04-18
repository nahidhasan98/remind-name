package feedback

import (
	"github.com/nahidhasan98/remind-name/config"
)

type repository struct {
	Collection string
}

func newRepository() *repository {
	return &repository{
		Collection: "feedback",
	}
}

func (repo *repository) saveFeedback(data *Feedback) error {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	_, err := collection.InsertOne(ctx, data)
	return err
}
