package feedback

import (
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	discordtexthook "github.com/nahidhasan98/discord-text-hook"
	"github.com/nahidhasan98/remind-name/config"
)

type handler struct {
	service *FeedbackService
}

func NewHandler() *handler {
	return &handler{
		service: NewFeedbackService(),
	}
}

func (h *handler) SaveFeedback(c *gin.Context) {
	var data Feedback

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trim, sanitize and escape potential dangerous characters in inputs
	data.Name = html.EscapeString(strings.TrimSpace(data.Name))
	data.Email = html.EscapeString(strings.TrimSpace(data.Email))
	data.Feedback = html.EscapeString(strings.TrimSpace(data.Feedback))

	// validate data
	if data.Name == "" || data.Email == "" || data.Feedback == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Email and Feedback are required fields."})
		return
	}

	res, err := h.service.SaveFeedback(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send to discord for instant notification
	go func() {
		disMsg := "```md\n"
		disMsg += "Name     : " + data.Name + "\n"
		disMsg += "Email    : " + data.Email + "\n"
		disMsg += "Feedback : " + data.Feedback + "\n"
		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(config.DISCORD_WEBHOOK_ID_FEEDBACK, config.DISCORD_WEBHOOK_TOKEN_FEEDBACK)
		ds.SendMessage(disMsg)
	}()

	c.JSON(http.StatusOK, res)
}
