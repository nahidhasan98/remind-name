package feedback

import (
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nahidhasan98/remind-name/helper"
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

	// Log the feedback
	log.Printf("New feedback received from %s (%s): %s", data.Name, data.Email, data.Feedback)

	// Send Discord notification
	helper.SendDiscordNotification("Feedback", data.Name, data.Email, data.Feedback)

	c.JSON(http.StatusOK, res)
}
