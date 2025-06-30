package notification

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nahidhasan98/remind-name/bot"
	"github.com/nahidhasan98/remind-name/logger"
	"github.com/nahidhasan98/remind-name/name"
	"github.com/nahidhasan98/remind-name/subscription"
)

// NotificationService handles notification related operations
type NotificationService struct {
	subscriptionService *subscription.SubscriptionService
	nameService         *name.NameService
	botManager          *bot.BotManager
}

// NewNotificationService creates a new instance of NotificationService
func NewNotificationService() *NotificationService {
	return &NotificationService{
		subscriptionService: subscription.NewSubscriptionService(),
		nameService:         name.NewNameService(),
		botManager:          bot.GetBotManager(),
	}
}

// generateMessage generates notification message for a subscription
func (ns *NotificationService) generateMessage(sub *subscription.Subscription) string {
	nextID := (sub.LastSentID % 99) + 1
	name, err := ns.nameService.GetName(nextID)
	if err != nil {
		logger.Error("Error fetching name: %v", err)
		return ""
	}

	// Using simpler bidirectional control:
	// \u200E (LRM) at start ensures left-to-right base direction
	// Individual LRM marks for each segment to maintain ordering
	return fmt.Sprintf("\u200E%d. %s \u200E\n[ %s ]: \u200E%s",
		name.ID,
		name.Languages["ar"].Transliteration,
		name.Languages["en"].Transliteration,
		name.Languages["en"].Meaning)
}

// sendNotification sends notification to a subscriber
func (ns *NotificationService) sendNotification(sub *subscription.Subscription) {
	message := ns.generateMessage(sub)
	if err := ns.botManager.SendMessage(sub, message); err != nil {
		logger.Error("Failed to send notification: %v", err)
	}
}

// checkAndSendNotifications checks subscriptions and sends notifications
func (ns *NotificationService) checkAndSendNotifications() {
	// last sent time is used in UTC seconds
	// GetSubscriptions() returns subscriptions that are active and eligible for notification regarding their last sent time and interval
	// time range and time zone are handled in this function
	currentUTCTime := time.Now().UTC().Unix()

	subcriptionService := subscription.NewSubscriptionService()
	subs, err := subcriptionService.GetSubscriptionsForDueNotification(currentUTCTime)
	if err != nil {
		logger.Error("Error fetching users: %v", err)
		return
	}

	const maxWorkers = 50
	workerCount := min(len(subs), maxWorkers)

	jobChan := make(chan subscription.Subscription)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for sub := range jobChan {
				loc, err := time.LoadLocation(sub.Timezone)
				if err != nil {
					logger.Error("Invalid timezone: %v %v", sub.Timezone, err)
					wg.Done()
					continue
				}

				now := time.Now().In(loc)
				currentTimeInSeconds := now.Hour()*3600 + now.Minute()*60

				if currentTimeInSeconds >= sub.TimeFrom && currentTimeInSeconds <= sub.TimeTo {
					ns.sendNotification(&sub)
					ns.subscriptionService.UpdateLastSent(&sub, currentUTCTime)
					logger.Info("Worker %d sent notification to %s", workerID, sub.Username)
				}
				wg.Done()
			}
		}(i)
	}

	// Send jobs
	for _, sub := range subs {
		wg.Add(1)
		jobChan <- sub
	}

	close(jobChan)
	wg.Wait() // Wait for all jobs to complete
}

// StartScheduler starts the notification schedule with context cancellation.
// This function runs in a goroutine and periodically sends notifications.
func (ns *NotificationService) StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	logger.Info("RUNNING: Notification scheduler.")

	// Run immediately on startup
	ns.checkAndSendNotifications()

	for {
		select {
		case <-ticker.C:
			ns.checkAndSendNotifications()
		case <-ctx.Done():
			logger.Info("Shutting down Notification scheduler...")
			return
		}
	}
}
