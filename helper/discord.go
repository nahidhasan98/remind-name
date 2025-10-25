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
		var id, token string

		disMsg := "```md\n"
		disMsg += fmt.Sprintf("# %s\n", title)

		switch title {
		case "Feedback":
			disMsg += "Name		: " + val1 + "\n"
			disMsg += "Email	: " + val2 + "\n"
			disMsg += "Feedback	: " + val3 + "\n"
			id, token = config.DISCORD_WEBHOOK_ID_FEEDBACK, config.DISCORD_WEBHOOK_TOKEN_FEEDBACK
		default:
			disMsg += "Platform	: " + val1 + "\n"
			disMsg += "Username	: " + val2 + "\n"
			disMsg += "Timezone	: " + val3 + "\n"
			id, token = config.DISCORD_WEBHOOK_ID_SUBSCRIPTION, config.DISCORD_WEBHOOK_TOKEN_SUBSCRIPTION
		}

		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(id, token)
		ds.SendMessage(disMsg)
		logger.Info("Sent Discord notification: %s, %s, %s", val1, val2, val3)
	}()
}
