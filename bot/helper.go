package bot

import (
	"errors"

	"github.com/nahidhasan98/remind-name/helper"
	"github.com/nahidhasan98/remind-name/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

func checkStatusAndGetMessage(username string, platform, user_idStr string) string {
	sub, err := subService.FindSubscription(username, platform, user_idStr)
	if err != nil {
		logger.Error("Error fetching subscription status: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) || sub == nil {
			return helper.NotSubscriber
		}
	}

	if sub != nil && sub.Status == 1 { // already verified subscription
		return `
Your subscription is already verified.

Here is your available options:
- Unsubscribe: Use /unsubscribe to stop receiving reminders.
- Help: Use /help to see all available commands.
`
	}

	return ""
}
