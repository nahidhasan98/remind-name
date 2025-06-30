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

func (repo *repository) getSubscriptionByUsernameAndPlatform(username string, platform string) (*Subscription, error) {
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

func (repo *repository) getSubscriptionForTelegramByUserID(user_idStr string) (*Subscription, error) {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	var sub Subscription
	err := collection.FindOne(ctx, bson.M{"username": user_idStr, "platform": "Telegram"}).Decode(&sub)
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
	filter := bson.M{"_id": data.ID}

	updateFields := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "schedule_type", Value: data.ScheduleType},
			{Key: "timezone", Value: data.Timezone},
			{Key: "time_from", Value: data.TimeFrom},
			{Key: "time_to", Value: data.TimeTo},
			{Key: "time_interval", Value: data.TimeInterval},
			{Key: "status", Value: data.Status},
			{Key: "token", Value: data.Token},
			{Key: "last_sent_at", Value: 0},
			{Key: "last_sent_id", Value: 0},
			{Key: "updated_at", Value: data.UpdatedAt},
		}},
	}

	return repo.updateSubscriptionFields(filter, updateFields)
}

func (repo *repository) updateSubscriptionStatus(sub *Subscription, status int8) error {
	filter := bson.M{"_id": sub.ID}

	updates := bson.D{{Key: "status", Value: status}, {Key: "updated_at", Value: time.Now().Unix()}}
	if sub.Platform == "Telegram" {
		updates = bson.D{
			{Key: "status", Value: status},
			{Key: "username", Value: sub.Username},
			{Key: "updated_at", Value: time.Now().Unix()},
		}
	}

	updateFields := bson.D{{Key: "$set", Value: updates}}
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

	filter := bson.M{"_id": sub.ID}
	updateFields := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_sent_at", Value: lastSentAt},
			{Key: "last_sent_id", Value: nextID},
			{Key: "updated_at", Value: time.Now().Unix()},
		}},
	}

	return repo.updateSubscriptionFields(filter, updateFields)
}

func (repo *repository) deleteUnverifiedSubscription(sub *Subscription) error {
	DB, ctx, cancel := config.DBConnect()
	defer cancel()
	defer DB.Client().Disconnect(ctx)

	collection := DB.Collection(repo.Collection)

	filter := bson.M{
		"username": sub.Username,
		"platform": sub.Platform,
		"status":   bson.M{"$in": []int8{0, 2}}, // Only delete unverified (0) or unsubscribed (2) subscriptions
	}

	_, err := collection.DeleteOne(ctx, filter)
	return err
}
