package helper

import (
	"fmt"

	"github.com/nahidhasan98/remind-name/logger"

	discordtexthook "github.com/nahidhasan98/discord-text-hook"
	"github.com/nahidhasan98/remind-name/config"
)

// SendDiscordNotification sends a Discord notification with the given title and subscription info
func SendDiscordNotification(title, val1, val2, val3 string) {
	go func() {
		disMsg := "```md\n"
		disMsg += fmt.Sprintf("# %s\n", title)

		switch title {
		case "Feedback":
			disMsg += "Name 	: " + val1 + "\n"
			disMsg += "Email 	: " + val2 + "\n"
			disMsg += "Feedback : " + val3 + "\n"
			disMsg += "```"
		default:
			disMsg += "Platform : " + val1 + "\n"
			disMsg += "Username : " + val2 + "\n"
			disMsg += "Timezone : " + val3 + "\n"
		}

		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(config.DISCORD_WEBHOOK_ID_SUBSCRIPTION, config.DISCORD_WEBHOOK_TOKEN_SUBSCRIPTION)
		ds.SendMessage(disMsg)
		logger.Info("Sent Discord notification: %s, %s, %s", val1, val2, val3)
	}()
}
