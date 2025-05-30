package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nahidhasan98/remind-name/config"
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

func checkStatusAndGetMessage(username, platform string) string {
	sub, err := subService.GetSubscription(username, platform)
	if err != nil {
		log.Printf("Error fetching subscription status: %v", err)
	}

	msg := ""

	if sub != nil && sub.Status == 1 { // already verified subscription
		msg = `
Your subscription already verified.

Here is your available options:
- Unsubscribe: Use /unsubscribe to stop receiving reminders.
- Help: Use /help to see all available commands.
`
	}

	return msg
}

func (t *TelegramBot) GetPlatformName() string {
	return t.Platform
}

// SendMessageToUser sends a message to a specific user
func (t *TelegramBot) SendMessageToUser(userID int64, message string) error {
	recipient := &tele.User{ID: userID}

	_, err := t.Bot.Send(recipient, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent to user %s: %s", recipient.Username, message)
	return nil
}

func (t *TelegramBot) Start(ctx context.Context) {
	bot := t.Bot

	// Handle /start command
	bot.Handle("/start", func(c tele.Context) error {
		msg := checkStatusAndGetMessage(c.Chat().Username, t.Platform)
		if msg != "" {
			return c.Send(msg)
		}

		welcomeMessage := `
Welcome to Remind Name Bot!

I can help you verify subscription and receive reminders.
Here's what you can do:
- Verify Subscription: Use /token <your-token> to verify your subscription.
- Unsubscribe: Use /unsubscribe to stop receiving reminders.
- Help: Use /help to see all available commands.

Let's get started!
`
		return c.Send(welcomeMessage)
	})

	// Handle /token command
	bot.Handle("/token", func(c tele.Context) error {
		msg := checkStatusAndGetMessage(c.Chat().Username, t.Platform)
		if msg != "" {
			return c.Send(msg)
		}

		args := c.Args()
		if len(args) != 1 {
			log.Printf("Invalid /token command by %s", c.Chat().Username)
			return c.Send("Usage: /token <your-token>")
		}

		token := args[0]

		err := subService.VerifySubscription(c.Chat().Username, t.Platform, token, c.Chat().ID)
		if err != nil {
			log.Printf("Failed verification for %s: %v", c.Chat().Username, err)
			return c.Send(err.Error())
		}

		log.Printf("User %s subscribed successfully.", c.Chat().Username)
		return c.Send("You have successfully subscribed!")
	})

	// Handle /unsubscribe command
	bot.Handle("/unsubscribe", func(c tele.Context) error {
		err := subService.Unsubscribe(c.Chat().Username, t.Platform)
		if err != nil {
			return c.Send(err.Error())
		}

		log.Printf("User %s unsubscribed successfully.", c.Chat().Username)
		return c.Send("You have successfully unsubscribed.")
	})

	// Add /help handler
	bot.Handle("/help", func(c tele.Context) error {
		helpMessage := `
Available Commands:

/start - Welcome message and instructions.
/token <your-token> - Verify your subscription.
/unsubscribe - Stop receiving reminders.
/help - Show this help message.

If you have any questions, feel free to reach out!
`
		return c.Send(helpMessage)
	})

	// Fallback handler for unknown commands
	bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Sorry, I didn't understand that command. Use /help to see the list of available commands.")
	})

	log.Println("RUNNING: Telegram Bot.")

	// Run the bot in a goroutine to handle shutdown
	go bot.Start()

	// Wait for the context to be canceled
	<-ctx.Done()
	log.Println("Shutting down Telegram bot...")

	// Remove the webhook
	dropPendingUpdates := true
	if err := bot.RemoveWebhook(dropPendingUpdates); err != nil {
		log.Printf("Error removing webhook: %v", err)
	}
	log.Println("Webhook removed successfully.")

	// Stop the bot gracefully
	// bot.Stop()
	log.Println("Telegram bot stopped.")
}
