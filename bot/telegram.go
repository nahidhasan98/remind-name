package bot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/nahidhasan98/remind-name/config"
	"github.com/nahidhasan98/remind-name/helper"
	"github.com/nahidhasan98/remind-name/logger"
	"github.com/nahidhasan98/remind-name/subscription"
	tele "gopkg.in/telebot.v4"
)

var subService = subscription.NewSubscriptionService()

var Telegram_Webhook = &tele.Webhook{
	Listen: fmt.Sprintf(":%v", config.APP_PORT),
	Endpoint: &tele.WebhookEndpoint{
		PublicURL: config.TELEGRAM_WEBHOOK_ENDPOINT,
	},
	DropUpdates: true,
}

var Telegram_LongPoller = &tele.LongPoller{
	Timeout: 10 * time.Second,
}

type TelegramBot struct {
	Platform string
	Bot      *tele.Bot
}

func NewTelegramBot() (*TelegramBot, error) {
	pref := tele.Settings{
		Token:  config.TELEGRAM_BOT_TOKEN_DEV,
		Poller: Telegram_LongPoller,
	}

	if config.APP_MODE == "production" {
		pref = tele.Settings{
			Token:  config.TELEGRAM_BOT_TOKEN,
			Poller: Telegram_Webhook,
		}
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bot: %v", err)
	}

	return &TelegramBot{
		Platform: "Telegram",
		Bot:      bot,
	}, nil
}

func (t *TelegramBot) GetPlatformName() string {
	return t.Platform
}

// SendMessageToUser sends a message to a specific user
func (t *TelegramBot) SendMessageToUser(username, message string) error {
	user_idStr, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid username: %v", err)
	}
	recipient := &tele.User{ID: user_idStr}

	_, err = t.Bot.Send(recipient, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	logger.Info("Message sent to user %s: %s", recipient.Username, message)
	return nil
}

func (t *TelegramBot) Start(ctx context.Context) {
	bot := t.Bot

	// Handle /start command
	bot.Handle("/start", func(c tele.Context) error {
		user_idStr := fmt.Sprintf("%d", c.Chat().ID)

		msg := checkStatusAndGetMessage(c.Chat().Username, t.Platform, user_idStr)
		if msg != "" && msg != helper.NotSubscriber {
			return c.Send(msg)
		}

		welcomeMessage := fmt.Sprintf(`
Assalamu alaikum %s
Your Username: %s
Your User ID: %d

Welcome to Remind Name Bot! I can help you to verify your subscription and receive notifications.

Here's what you can do:
- Verify Subscription: (After subscribing from website) Use /token <your-token> to verify your subscription.
- Unsubscribe: Use /unsubscribe to stop receiving notifications.
- Help: Use /help to see all available commands.

Let's get started!
`, c.Chat().FirstName, c.Chat().Username, c.Chat().ID)
		return c.Send(welcomeMessage)
	})

	// Handle /token command
	bot.Handle("/token", func(c tele.Context) error {
		args := c.Args()
		if len(args) != 1 {
			logger.Warn("Invalid /token command by %s", c.Chat().Username)
			return c.Send("Usage: /token <your-token>")
		}

		token := args[0]
		user_idStr := fmt.Sprintf("%d", c.Chat().ID)

		err := subService.VerifySubscription(c.Chat().Username, t.Platform, token, user_idStr)
		if err != nil {
			logger.Error("Failed verification for %s: %v", c.Chat().Username, err)
			return c.Send(helper.FormatErrorMessage(err.Error()))
		}

		logger.Info("User %s subscribed successfully.", c.Chat().Username)
		return c.Send("You have successfully subscribed!")
	})

	// Handle /unsubscribe command
	bot.Handle("/unsubscribe", func(c tele.Context) error {
		id_str := fmt.Sprintf("%d", c.Chat().ID)
		err := subService.Unsubscribe(c.Chat().Username, t.Platform, id_str)
		if err != nil {
			return c.Send(helper.FormatErrorMessage(err.Error()))
		}

		logger.Info("User %s unsubscribed successfully.", c.Chat().Username)
		return c.Send("You have successfully unsubscribed.")
	})

	// Add /help handler
	bot.Handle("/help", func(c tele.Context) error {
		helpMessage := `
Available Commands:

/start - Welcome message and instructions.
/token <your-token> - Verify your subscription.
/unsubscribe - Stop receiving notifications.
/help - Show this help message.

If you have any questions, feel free to reach out!
https://remind.name
`
		return c.Send(helpMessage)
	})

	// Fallback handler for unknown commands
	bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Sorry, I didn't understand that command. Use /help to see the list of available commands.")
	})

	logger.Info("RUNNING: Telegram Bot.")

	// Run the bot in a goroutine to handle shutdown
	go bot.Start()

	// Wait for the context to be canceled
	<-ctx.Done()
	logger.Info("Shutting down Telegram bot...")

	// Remove the webhook
	dropPendingUpdates := true
	if err := bot.RemoveWebhook(dropPendingUpdates); err != nil {
		logger.Error("Error removing webhook: %v", err)
	}
	logger.Info("Webhook removed successfully.")

	// Stop the bot gracefully
	// bot.Stop()
	logger.Info("Telegram bot stopped.")
}
