package subscription

import (
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	discordtexthook "github.com/nahidhasan98/discord-text-hook"
	"github.com/nahidhasan98/remind-name/config"
)

type handler struct {
	service *SubscriptionService
}

func NewHandler() *handler {
	return &handler{
		service: NewSubscriptionService(),
	}
}

// Helper function to validate platform inputs
func isValidPlatform(name string) bool {
	validPlatforms, err := platformService.GetAllPlatforms()
	if err != nil {
		return false
	}

	for _, p := range validPlatforms {
		if p.Name == name {
			return true
		}
	}

	return false
}

// Helper function to validate time inputs
func isValidTime(hours, minutes, ampm string) (bool, int) {
	h, err := strconv.Atoi(hours)
	if ampm == "ignore" && (err != nil || h < 1 || h > 23) {
		return false, 0
	} else if err != nil || h < 1 || h > 12 {
		return false, 0
	}

	m, err := strconv.Atoi(minutes)
	if err != nil || m < 0 || m > 59 {
		return false, 0
	}

	if ampm != "am" && ampm != "pm" && ampm != "ignore" {
		return false, 0
	}

	// Convert time to seconds
	if ampm == "pm" && h != 12 {
		h += 12
	} else if ampm == "am" && h == 12 {
		h = 0
	}
	seconds := (h * 60 * 60) + (m * 60)

	return true, seconds
}

// Helper function to validate timezone inputs
func validateTimezone(tz string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "UTC"
	}

	return loc.String()
}

func (h *handler) AddSubscription(c *gin.Context) {
	var sub Subscription

	if err := c.ShouldBind(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trim, sanitize and escape potential dangerous characters in inputs
	sub.Platform = html.EscapeString(strings.TrimSpace(sub.Platform))
	sub.Username = html.EscapeString(strings.TrimSpace(sub.Username))
	sub.ScheduleType = html.EscapeString(strings.TrimSpace(sub.ScheduleType))
	sub.FromHour = html.EscapeString(strings.TrimSpace(c.PostForm("fromHour")))
	// sub.FromMinute = html.EscapeString(strings.TrimSpace(c.PostForm("fromMinute")))
	sub.FromMinute = "00"
	sub.FromAMPM = html.EscapeString(strings.TrimSpace(c.PostForm("fromAMPM")))
	sub.ToHour = html.EscapeString(strings.TrimSpace(c.PostForm("toHour")))
	// sub.ToMinute = html.EscapeString(strings.TrimSpace(c.PostForm("toMinute")))
	sub.ToMinute = "00"
	sub.ToAMPM = html.EscapeString(strings.TrimSpace(c.PostForm("toAMPM")))
	sub.IntervalHour = html.EscapeString(strings.TrimSpace(c.PostForm("intervalHour")))
	// sub.IntervalMinute = html.EscapeString(strings.TrimSpace(c.PostForm("intervalMinute")))
	sub.IntervalMinute = "00"
	sub.Timezone = html.EscapeString(strings.TrimSpace(sub.Timezone))

	// Validate data
	if sub.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	sub.Timezone = validateTimezone(sub.Timezone)

	// Validate platform and schedule values
	if !isValidPlatform(sub.Platform) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid platform"})
		return
	}

	// Validate schedule
	if sub.ScheduleType != "default" && sub.ScheduleType != "custom" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule type"})
		return
	}

	if sub.ScheduleType == "default" {
		sub.TimeFrom = 6 * 60 * 60     // 6:00 AM
		sub.TimeTo = 21 * 60 * 60      // 9:00 PM
		sub.TimeInterval = 3 * 60 * 60 // 3:00 Hour
	} else {
		valid, fromSeconds := isValidTime(sub.FromHour, sub.FromMinute, sub.FromAMPM)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' time"})
			return
		}

		valid, toSeconds := isValidTime(sub.ToHour, sub.ToMinute, sub.ToAMPM)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' time"})
			return
		}

		valid, intervalSeconds := isValidTime(sub.IntervalHour, sub.IntervalMinute, "ignore")
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'interval' time"})
			return
		}

		sub.TimeFrom = fromSeconds
		sub.TimeTo = toSeconds
		sub.TimeInterval = intervalSeconds

		if sub.TimeFrom > sub.TimeTo {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timing.<br>'From' time should be less than or equal to 'To' time."})
			return
		}

		if sub.TimeInterval < 30*60 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timing.<br>'Interval' should be at least 30 minutes."})
			return
		}
	}

	res, err := h.service.AddSubscription(&sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// send to discord for instant notification
	go func() {
		disMsg := "```md\n"
		disMsg += "# Subscribed\n"
		disMsg += "Platform : " + sub.Platform + "\n"
		disMsg += "Username : " + sub.Username + "\n"
		disMsg += "Timezone : " + sub.Timezone + "\n"
		disMsg += "```"

		ds := discordtexthook.NewDiscordTextHookService(config.DISCORD_WEBHOOK_ID_SUBSCRIPTION, config.DISCORD_WEBHOOK_TOKEN_SUBSCRIPTION)
		ds.SendMessage(disMsg)
	}()

	c.JSON(http.StatusOK, res)
}
