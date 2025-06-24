package bot

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nahidhasan98/remind-name/subscription"
)

// Bot is implemented by each platform-specific bot.
type Bot interface {
	Start(ctx context.Context)
	SendMessageToUser(username, message string) error
	GetPlatformName() string
}

// BotManager manages a set of bots.
type BotManager struct {
	bots map[string]Bot
}

var (
	managerInstance *BotManager
	once            sync.Once
)

// GetBotManager returns the singleton instance of BotManager.
// It initializes all supported bots once.
func GetBotManager() *BotManager {
	once.Do(func() {
		bots := make(map[string]Bot)

		// Example: initialize Telegram bot once.
		telegramBot, err := NewTelegramBot()
		if err != nil {
			log.Printf("Error initializing Telegram bot: %v", err)
		} else {
			bots["Telegram"] = telegramBot
		}

		// (Optional) Initialize other bots (Discord, WhatsApp, etc.)
		// discordBot, err := NewDiscordBot()
		// if err != nil {
		// 	log.Printf("Error initializing Discord bot: %v", err)
		// } else {
		// 	bots["Discord"] = discordBot
		// }

		managerInstance = &BotManager{
			bots: bots,
		}
	})
	return managerInstance
}

// StartAll calls the Start method on each bot concurrently.
func (m *BotManager) StartAll(ctx context.Context) {
	for _, bot := range m.bots {
		// go bot.Start(ctx)
		bot.Start(ctx)
	}
}

// SendMessage looks up the bot by the subscription's Platform field
// and sends the message to the user.
func (m *BotManager) SendMessage(sub *subscription.Subscription, message string) error {
	bot, exists := m.bots[sub.Platform]
	if !exists {
		return fmt.Errorf("platform %q is not supported", sub.Platform)
	}

	if err := bot.SendMessageToUser(sub.Username, message); err != nil {
		return fmt.Errorf("error sending message via %s: %w", sub.Platform, err)
	}
	return nil
}
