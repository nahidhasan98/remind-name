package helper

import (
	"fmt"

	discordtexthook "github.com/nahidhasan98/discord-text-hook"
	"github.com/nahidhasan98/remind-name/config"
)

// SendDiscordNotification sends a Discord notification with the given title and subscription info
func SendDiscordNotification(title, platform, username, timezone string) {
	go func() {
		disMsg := "```md\n"
		disMsg += fmt.Sprintf("# %s\n", title)
		disMsg += "Platform : " + platform + "\n"
		disMsg += "Username : " + username + "\n"
		disMsg += "Timezone : " + timezone + "\n"
		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(config.DISCORD_WEBHOOK_ID_SUBSCRIPTION, config.DISCORD_WEBHOOK_TOKEN_SUBSCRIPTION)
		ds.SendMessage(disMsg)
	}()
}
