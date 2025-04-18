package subscription

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	discordtexthook "github.com/nahidhasan98/discord-text-hook"
	"github.com/nahidhasan98/remind-name/config"
	"github.com/nahidhasan98/remind-name/platform"
	"go.mongodb.org/mongo-driver/mongo"
)

var platformService = platform.NewPlatformService()

type SubscriptionService struct {
	repo *repository
}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{
		repo: newRepository(),
	}
}

// generateToken generates a random token of the specified length
func generateToken(length int) (string, error) {
	// Create a byte slice of the desired length
	token := make([]byte, length)

	// Fill the slice with random bytes
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	/*
		Using RawURLEncoding omits the padding characters like (==)
		Calculation: Base64 length = ceil(length/3) * 4
		So length 16 gives 22 characters with padding
	*/

	// Encode the byte slice into a URL-safe base64 string
	return base64.RawURLEncoding.EncodeToString(token), nil
}

func (service *SubscriptionService) AddSubscription(data *Subscription) (*Response, error) {
	sub, err := service.repo.getSubscriptionByUsernameAndPlatform(data.Username, data.Platform)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("Error fetching subscription for username: %s, platform: %s, error: %v", data.Username, data.Platform, err)
		return nil, errors.New("Something went wrong. Please try again.")
	}

	platform, err := platformService.GetPlatformDetailsByName(data.Platform)
	if err != nil {
		log.Printf("Error fetching platform details for, platform: %s, error: %v", data.Platform, err)
		return nil, errors.New("Something went wrong. Please try again.")
	}

	// Handle existing subscription cases
	if sub != nil {
		if sub.Status == 1 {
			return &Response{
				Status:  sub.Status,
				Message: "You are already a verified subscriber.",
			}, nil
		}

		if sub.Status == 0 {
			return &Response{
				Platform: platform,
				Token:    sub.Token,
				Status:   sub.Status,
				Message:  fmt.Sprintf("You have already subscribed.\nPlease verify your %s ID.", data.Platform),
			}, nil
		}
	}

	// Handle new subscription or unsubscribed case
	token, err := generateToken(16)
	if err != nil {
		log.Printf("Error generating token for username: %s, platform: %s, error: %v", data.Username, data.Platform, err)
		return nil, errors.New("Something went wrong. Please try again.")
	}

	data.Token = token
	data.Status = 0
	now := time.Now().Unix()
	data.UpdatedAt = now
	data.CreatedAt = now

	if sub != nil && sub.Status == 2 {
		err = service.repo.updateSubscription(data)
	} else {
		err = service.repo.createSubscription(data)
	}

	if err != nil {
		return nil, err
	}

	return &Response{
		Platform: platform,
		Token:    data.Token,
		Status:   data.Status,
		Message:  "Subscribed successfully!",
	}, nil
}

func (service *SubscriptionService) VerifySubscription(username, platform, token string, userID int64) error {
	sub, err := service.repo.getSubscriptionByUsernameAndPlatform(username, platform)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("You are not a subscriber. Please subscribe first from\nhttps://remind.name")
		}
		log.Printf("Error fetching subscription for username: %s, platform: %s, error: %v", username, platform, err)
		return errors.New("Something went wrong. Please try again.")
	}

	if sub.Status == 2 {
		return errors.New("You are not a subscriber. Please subscribe first from\nhttps://remind.name")
	}

	if sub.Status == 1 {
		return errors.New("Your subscription already verified.")
	}

	if sub.Token != token {
		return errors.New("Invalid token. Please try again.")
	}

	if err := service.repo.updateSubscriptionToVerified(username, platform, userID); err != nil {
		log.Printf("Error updating subscription status for username: %s, platform: %s, error: %v", username, platform, err)
		return errors.New("Something went wrong. Please try again.")
	}

	// send to discord for instant notification
	go func() {
		disMsg := "```md\n"
		disMsg += "# Verified\n"
		disMsg += "Platform : " + sub.Platform + "\n"
		disMsg += "Username : " + sub.Username + "\n"
		disMsg += "Timezone : " + sub.Timezone + "\n"
		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(config.DISCORD_WEBHOOK_ID_SUBSCRIPTION, config.DISCORD_WEBHOOK_TOKEN_SUBSCRIPTION)
		ds.SendMessage(disMsg)
	}()

	return nil
}

func (service *SubscriptionService) Unsubscribe(username, platform string) error {
	sub, err := service.repo.getSubscriptionByUsernameAndPlatform(username, platform)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("You are not a subscriber. Please subscribe first from\nhttps://remind.name")
		}
		log.Printf("Error fetching subscription for username: %s, platform: %s, error: %v", username, platform, err)
		return errors.New("Something went wrong. Please try again.")
	}

	if sub.Status == 2 {
		return errors.New("You are not a subscriber. Please subscribe first from\nhttps://remind.name")
	}

	if err := service.repo.updateSubscriptionToUnsubscribed(username, platform); err != nil {
		log.Printf("Error updating subscription status for username: %s, platform: %s, error: %v", username, platform, err)
		return errors.New("Something went wrong. Please try again.")
	}

	// send to discord for instant notification
	go func() {
		disMsg := "```md\n"
		disMsg += "# Unsubscribed\n"
		disMsg += "Platform : " + sub.Platform + "\n"
		disMsg += "Username : " + sub.Username + "\n"
		disMsg += "Timezone : " + sub.Timezone + "\n"
		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(config.DISCORD_WEBHOOK_ID_SUBSCRIPTION, config.DISCORD_WEBHOOK_TOKEN_SUBSCRIPTION)
		ds.SendMessage(disMsg)
	}()

	return nil
}

func (service *SubscriptionService) GetSubscription(username, platform string) (*Subscription, error) {
	sub, err := service.repo.getSubscriptionByUsernameAndPlatform(username, platform)
	if err != nil {
		log.Printf("Error fetching subscription for username: %s, platform: %s, error: %v", username, platform, err)
		return nil, err
	}

	return sub, nil
}

func (service *SubscriptionService) GetSubscriptionsForDueNotification(currentUTCTime int64) ([]Subscription, error) {
	subs, err := service.repo.getSubscriptionsForDueNotification(currentUTCTime)
	if err != nil {
		log.Printf("Error fetching subscriptions, error: %v", err)
		return nil, errors.New("Something went wrong. Please try again.")
	}

	return subs, nil
}

func (service *SubscriptionService) UpdateLastSent(sub *Subscription, lastSentAt int64) error {
	if err := service.repo.updateLastSent(sub, lastSentAt); err != nil {
		log.Printf("Error updating last sent for username: %s, error: %v", sub.Username, err)
		return errors.New("Something went wrong. Please try again.")
	}

	return nil
}
