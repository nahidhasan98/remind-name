package subscription

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/nahidhasan98/remind-name/helper"
	"github.com/nahidhasan98/remind-name/logger"
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
		logger.Error("Error fetching subscription for username: %s, platform: %s, error: %v", data.Username, data.Platform, err)
		return nil, errors.New(helper.SomethingWentWrong)
	}

	platform, err := platformService.GetPlatformDetailsByName(data.Platform)
	if err != nil {
		logger.Error("Error fetching platform details for, platform: %s, error: %v", data.Platform, err)
		return nil, errors.New(helper.SomethingWentWrong)
	}

	// Handle existing subscription cases
	if sub != nil {
		if sub.Status == 1 {
			return &Response{
				Status:  sub.Status,
				Message: helper.FormatErrorMessage(helper.AlreadyVerifiedSub),
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
		logger.Error("Error generating token for username: %s, platform: %s, error: %v", data.Username, data.Platform, err)
		return nil, errors.New(helper.SomethingWentWrong)
	}

	now := time.Now().Unix()

	if sub != nil && sub.Status == 2 {
		sub.Token = token
		sub.Status = 0
		sub.ScheduleType = data.ScheduleType
		sub.Timezone = data.Timezone
		sub.TimeFrom = data.TimeFrom
		sub.TimeTo = data.TimeTo
		sub.TimeInterval = data.TimeInterval
		sub.LastSentAt = 0
		sub.LastSentID = 0
		sub.UpdatedAt = now
		err = service.repo.updateSubscription(sub)
	} else {
		data.Token = token
		data.Status = 0
		data.UpdatedAt = now
		data.CreatedAt = now
		err = service.repo.createSubscription(data)
	}

	if err != nil {
		return nil, err
	}

	return &Response{
		Platform: platform,
		Token:    token,
		Status:   0,
		Message:  "Subscribed successfully!",
	}, nil
}

// validateSubscriptionStatus checks the subscription status and token, returns error if invalid
func validateSubscriptionStatus(sub *Subscription, token string) error {
	if sub.Status == 2 {
		return helper.ErrNotSubscriber
	}
	if sub.Status == 1 {
		return helper.ErrAlreadyVerified
	}
	if sub.Token != token {
		return helper.ErrInvalidToken
	}
	return nil
}

func (service *SubscriptionService) VerifySubscription(username, platform, token, user_idStr string) error {
	var sub *Subscription
	var err error

	if platform != "Telegram" {
		sub, err = service.repo.getSubscriptionByUsernameAndPlatform(username, platform)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return helper.ErrNotSubscriber
			}
			logger.Error("Error fetching subscription for username: %s, platform: %s, error: %v", username, platform, err)
			return helper.ErrSomethingWentWrong
		}

		if err := validateSubscriptionStatus(sub, token); err != nil {
			return err
		}
	} else {
		// For Telegram platform, try by user_idStr first, then by username if needed
		var foundByUserID bool = false
		// Try by user_idStr first
		sub, err = service.repo.getSubscriptionForTelegramByUserID(user_idStr)
		if err == nil && sub != nil {
			// Check status and token
			if err := validateSubscriptionStatus(sub, token); err == nil {
				foundByUserID = true
			} else if errors.Is(err, helper.ErrAlreadyVerified) {
				return err
			} else {
				foundByUserID = false
			}
		}
		unverifiedSub := &Subscription{Username: username, Platform: "Telegram"}

		if !foundByUserID {
			// Try by username
			sub, err = service.repo.getSubscriptionByUsernameAndPlatform(username, "Telegram")
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return helper.ErrNotSubscriber
				}
				logger.Error("Error fetching subscription for username: %s, platform: %s, error: %v", username, platform, err)
				return helper.ErrSomethingWentWrong
			}
			if err := validateSubscriptionStatus(sub, token); err != nil {
				return err
			}

			// Update username to user_idStr for Telegram
			sub.Username = user_idStr

			// Delete unverified subscription for username if exists
			unverifiedSub = &Subscription{Username: user_idStr, Platform: "Telegram"}
		}
		err = service.repo.deleteUnverifiedSubscription(unverifiedSub)
		if err != nil {
			logger.Error("Error deleting unverified subscription for username: %s, platform: %s, error: %v", username, platform, err)
			return helper.ErrSomethingWentWrong
		}
		// log deleteing info
		logger.Info("Deleted unverified subscription for username: %s, platform: %s", unverifiedSub.Username, unverifiedSub.Platform)
	}

	// If found by user_idStr and all checks pass, update status and delete unverified subscription for username
	if err := service.repo.updateSubscriptionStatus(sub, 1); err != nil {
		logger.Error("Error updating subscription status for username: %s, platform: %s, error: %v", username, platform, err)
		return helper.ErrSomethingWentWrong
	}

	// Send Discord notification
	helper.SendDiscordNotification("Verified", sub.Platform, sub.Username, sub.Timezone)
	return nil
}

func (service *SubscriptionService) Unsubscribe(username, platform, user_idStr string) error {
	var sub *Subscription
	var err error

	if platform != "Telegram" {
		sub, err = service.repo.getSubscriptionByUsernameAndPlatform(username, platform)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return helper.ErrNotSubscriber
			}
			logger.Error("Error fetching subscription for username: %s, platform: %s, error: %v", username, platform, err)
			return helper.ErrSomethingWentWrong
		}

		if sub.Status == 2 {
			return errors.New(helper.NotSubscriber)
		}
		if sub.Status == 0 {
			return errors.New(helper.UnverifiedSubscriber)
		}
		// Mark as unsubscribed
		if err := service.repo.updateSubscriptionStatus(sub, 2); err != nil {
			logger.Error("Error updating subscription status for username: %s, platform: %s, error: %v", username, platform, err)
			return errors.New(helper.SomethingWentWrong)
		}
		// Send Discord notification
		helper.SendDiscordNotification("Unsubscribed", sub.Platform, sub.Username, sub.Timezone)
		return nil
	}

	// Telegram logic
	subByID, _ := service.repo.getSubscriptionForTelegramByUserID(user_idStr)
	subByUser, _ := service.repo.getSubscriptionByUsernameAndPlatform(username, "Telegram")

	// If subByID is nil, check subByUser
	if subByID == nil {
		if subByUser != nil && subByUser.Status == 0 {
			logger.Warn("Unverified subscription for username: %s, platform: %s", username, platform)
			return errors.New(helper.UnverifiedSubscriber)
		}
		logger.Warn("No subscription found for username: %s, platform: %s", username, platform)
		return helper.ErrNotSubscriber
	}

	// subByID is not nil
	if subByID.Status == 0 {
		logger.Warn("Unverified subscription for user_idStr: %s, platform: %s", user_idStr, platform)
		return errors.New(helper.UnverifiedSubscriber)
	}
	if subByID.Status == 2 {
		logger.Warn("Subscription already unsubscribed for user_idStr: %s, platform: %s", user_idStr, platform)
		return errors.New(helper.NotSubscriber)
	}

	// subByID.Status == 1, proceed to unsubscribe
	if err := service.repo.updateSubscriptionStatus(subByID, 2); err != nil {
		logger.Error("Error updating subscription status for user_idStr: %s, platform: %s, error: %v", user_idStr, platform, err)
		return errors.New(helper.SomethingWentWrong)
	}

	// Delete unverified subByUser if exists
	if subByUser != nil && subByUser.Status == 0 {
		toDelete := &Subscription{Username: username, Platform: "Telegram"}
		err = service.repo.deleteUnverifiedSubscription(toDelete)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error("Error deleting unverified subscription for username: %s, platform: %s, error: %v", username, platform, err)
			return helper.ErrSomethingWentWrong
		}
		logger.Info("Deleted unverified subscription for username: %s, platform: %s", toDelete.Username, toDelete.Platform)
	}

	helper.SendDiscordNotification("Unsubscribed", subByID.Platform, subByID.Username, subByID.Timezone)
	return nil
}

func (service *SubscriptionService) GetSubscriptionsForDueNotification(currentUTCTime int64) ([]Subscription, error) {
	subs, err := service.repo.getSubscriptionsForDueNotification(currentUTCTime)
	if err != nil {
		logger.Error("Error fetching subscriptions, error: %v", err)
		return nil, errors.New(helper.SomethingWentWrong)
	}

	return subs, nil
}

func (service *SubscriptionService) UpdateLastSent(sub *Subscription, lastSentAt int64) error {
	if err := service.repo.updateLastSent(sub, lastSentAt); err != nil {
		logger.Error("Error updating last sent for username: %s, error: %v", sub.Username, err)
		return errors.New(helper.SomethingWentWrong)
	}

	return nil
}

// FindSubscription finds a subscription by username/platform/user_idStr, handling Telegram logic
func (service *SubscriptionService) FindSubscription(username, platform, user_idStr string) (*Subscription, error) {
	if platform == "Telegram" {
		sub, err := service.repo.getSubscriptionForTelegramByUserID(user_idStr)
		if (err == nil && sub != nil) || (err != nil && err != mongo.ErrNoDocuments) {
			return sub, err
		}
		// Try by username if not found by user_idStr
		return service.repo.getSubscriptionByUsernameAndPlatform(username, "Telegram")
	}
	return service.repo.getSubscriptionByUsernameAndPlatform(username, platform)
}
