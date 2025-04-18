package subscription

import (
	"context"
	"time"

	"github.com/nahidhasan98/remind-name/config"
	"go.mongodb.org/mongo-driver/bson"
)

type repository struct {
	Collection string
}

func newRepository() *repository {
	return &repository{
		Collection: "subscription",
	}
}

func (repo *repository) getSubscriptionByUsernameAndPlatform(username, platform string) (*Subscription, error) {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	var sub Subscription
	err := collection.FindOne(ctx, bson.M{"username": username, "platform": platform}).Decode(&sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (repo *repository) createSubscription(data *Subscription) error {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	_, err := collection.InsertOne(ctx, data)
	return err
}

func (repo *repository) updateSubscriptionFields(filter bson.M, updateFields bson.D) error {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	_, err := collection.UpdateOne(ctx, filter, updateFields)
	return err
}

func (repo *repository) updateSubscription(data *Subscription) error {
	filter := bson.M{"username": data.Username, "platform": data.Platform}

	updateFields := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "schedule_type", Value: data.ScheduleType},
			{Key: "time_from", Value: data.TimeFrom},
			{Key: "time_to", Value: data.TimeTo},
			{Key: "time_interval", Value: data.TimeInterval},
			{Key: "status", Value: data.Status},
			{Key: "updated_at", Value: data.UpdatedAt},
		}},
	}

	return repo.updateSubscriptionFields(filter, updateFields)
}

func (repo *repository) updateSubscriptionToVerified(username, platform string, userID int64) error {
	filter := bson.M{"username": username, "platform": platform}

	updateFields := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "status", Value: 1},
			{Key: "updated_at", Value: time.Now().Unix()},
		}},
	}

	return repo.updateSubscriptionFields(filter, updateFields)
}

func (repo *repository) updateSubscriptionToUnsubscribed(username, platform string) error {
	filter := bson.M{"username": username, "platform": platform}

	updateFields := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: 2},
			{Key: "updated_at", Value: time.Now().Unix()},
		}},
	}

	return repo.updateSubscriptionFields(filter, updateFields)
}

func (repo *repository) getSubscriptionsForDueNotification(currentUTCTime int64) ([]Subscription, error) {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	filter := bson.M{
		"$and": []bson.M{
			{
				"$expr": bson.M{
					"$gte": []interface{}{
						bson.M{"$subtract": []interface{}{currentUTCTime, "$last_sent_at"}},
						"$time_interval",
					},
				},
			},
			{"status": 1}, // Only active subscriptions
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var subs []Subscription
	if err = cursor.All(context.TODO(), &subs); err != nil {
		return nil, err
	}

	return subs, nil
}

func (repo *repository) updateLastSent(sub *Subscription, lastSentAt int64) error {
	nextID := (sub.LastSentID % 99) + 1

	filter := bson.M{"username": sub.Username}
	updateFields := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_sent_at", Value: lastSentAt},
			{Key: "last_sent_id", Value: nextID},
			{Key: "updated_at", Value: time.Now().Unix()},
		}},
	}

	return repo.updateSubscriptionFields(filter, updateFields)
}
