package subscription

import (
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nahidhasan98/remind-name/helper"
	"github.com/nahidhasan98/remind-name/logger"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(err.Error())})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.UsernameRequired)})
		return
	}

	if sub.Username[0] == '@' {
		sub.Username = sub.Username[1:]
	}

	sub.Timezone = validateTimezone(sub.Timezone)

	// Validate platform and schedule values
	if !isValidPlatform(sub.Platform) {
		c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidPlatform)})
		return
	}

	// Validate schedule
	if sub.ScheduleType != "default" && sub.ScheduleType != "custom" {
		c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidScheduleType)})
		return
	}

	if sub.ScheduleType == "default" {
		sub.TimeFrom = 6 * 60 * 60     // 6:00 AM
		sub.TimeTo = 21 * 60 * 60      // 9:00 PM
		sub.TimeInterval = 3 * 60 * 60 // 3:00 Hour
	} else {
		okFrom, fromSeconds := isValidTime(sub.FromHour, sub.FromMinute, sub.FromAMPM)
		if !okFrom {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidFromTime)})
			return
		}

		okTo, toSeconds := isValidTime(sub.ToHour, sub.ToMinute, sub.ToAMPM)
		if !okTo {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidToTime)})
			return
		}

		okInterval, intervalSeconds := isValidTime(sub.IntervalHour, sub.IntervalMinute, "ignore")
		if !okInterval {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidIntervalTime)})
			return
		}

		sub.TimeFrom = fromSeconds
		sub.TimeTo = toSeconds
		sub.TimeInterval = intervalSeconds

		if sub.TimeFrom > sub.TimeTo {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidTimingFromTo)})
			return
		}

		if sub.TimeInterval < 30*60 {
			c.JSON(http.StatusBadRequest, gin.H{"error": helper.FormatErrorMessage(helper.InvalidTimingInterval)})
			return
		}
	}

	res, err := h.service.AddSubscription(&sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": helper.FormatErrorMessage(err.Error())})
		return
	}

	// Log the subscription
	logger.Info("New subscription added: %s (%s) - %s", sub.Username, sub.Platform, sub.Timezone)

	// Send Discord notification
	helper.SendDiscordNotification("Subscribed", sub.Platform, sub.Username, sub.Timezone)

	c.JSON(http.StatusOK, res)
}
