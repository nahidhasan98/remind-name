package subscription

import "github.com/nahidhasan98/remind-name/platform"

// Subscription represents a user's subscription to receive reminders.
type Subscription struct {
	ID           interface{} `bson:"_id,omitempty"`
	Platform     string      `bson:"platform" form:"platform" binding:"required"`
	Username     string      `bson:"username" form:"username" binding:"required"`
	ScheduleType string      `bson:"schedule_type" form:"scheduleType" binding:"required"`
	Timezone     string      `bson:"timezone" form:"timezone"`
	TimeFrom     int         `bson:"time_from"`
	TimeTo       int         `bson:"time_to"`
	TimeInterval int         `bson:"time_interval"`
	Status       int8        `bson:"status"` //  status: 0: not verified, 1: verified, 2: unsubscribed
	LastSentAt   int64       `bson:"last_sent_at"`
	LastSentID   int         `bson:"last_sent_id"`
	Token        string      `bson:"token"`
	CreatedAt    int64       `bson:"created_at"`
	UpdatedAt    int64       `bson:"updated_at"`

	// Fields for form only (not stored in DB)
	FromHour       string `bson:"-" form:"fromHour"`
	FromMinute     string `bson:"-" form:"fromMinute"`
	FromAMPM       string `bson:"-" form:"fromAMPM"`
	ToHour         string `bson:"-" form:"toHour"`
	ToMinute       string `bson:"-" form:"toMinute"`
	ToAMPM         string `bson:"-" form:"toAMPM"`
	IntervalHour   string `bson:"-" form:"intervalHour"`
	IntervalMinute string `bson:"-" form:"intervalMinute"`
}

type Response struct {
	Platform *platform.Platform `json:",omitempty"`
	Token    string             `json:",omitempty"`
	Status   int8
	Message  string
}
