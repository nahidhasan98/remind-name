package helper

import "errors"

const (
	NotSubscriber        = "not_subscriber"
	AlreadyVerified      = "already_verified"
	InvalidToken         = "invalid_token"
	UnverifiedSubscriber = "unverified_subscriber"
	SomethingWentWrong   = "something_went_wrong"
	AlreadyVerifiedSub   = "already_verified_sub"

	UsernameRequired      = "username is required"
	InvalidPlatform       = "invalid platform"
	InvalidScheduleType   = "invalid schedule type"
	InvalidFromTime       = "invalid 'from' time"
	InvalidToTime         = "invalid 'to' time"
	InvalidIntervalTime   = "invalid 'interval' time"
	InvalidTimingFromTo   = "invalid timing.<br>'From' time should be less than or equal to 'To' time."
	InvalidTimingInterval = "invalid timing. 'interval' should be at least 30 minutes."
)

var ErrorDisplayMessages = map[string]string{
	NotSubscriber:         "You are not a subscriber. Please subscribe first from\nhttps://remind.name",
	AlreadyVerified:       "Your subscription is already verified.",
	InvalidToken:          "Invalid token. Please try again.",
	UnverifiedSubscriber:  "You are not a verified subscriber. Nothing to unsubscribe.",
	SomethingWentWrong:    "Something went wrong. Please try again.",
	AlreadyVerifiedSub:    "You are already a verified subscriber.",
	UsernameRequired:      "Username is required.",
	InvalidPlatform:       "Invalid platform.",
	InvalidScheduleType:   "Invalid schedule type.",
	InvalidFromTime:       "Invalid 'from' time.",
	InvalidToTime:         "Invalid 'to' time.",
	InvalidIntervalTime:   "Invalid 'interval' time.",
	InvalidTimingFromTo:   "Invalid timing. 'From' time should be less than or equal to 'To' time.",
	InvalidTimingInterval: "Invalid timing. 'Interval' should be at least 30 minutes.",
}

// FormatErrorMessage returns a user-facing error message for a given error key
func FormatErrorMessage(msg string) string {
	if display, ok := ErrorDisplayMessages[msg]; ok {
		return display
	}
	// fallback: capitalize and punctuate
	if len(msg) == 0 {
		return msg
	}
	// Capitalize first letter
	formatted := string(msg[0]-32) + msg[1:]
	// Add period if not present
	if formatted[len(formatted)-1] != '.' {
		formatted += "."
	}
	return formatted
}

// Error variables for use with errors.Is
var (
	ErrNotSubscriber        = errors.New(NotSubscriber)
	ErrAlreadyVerified      = errors.New(AlreadyVerified)
	ErrInvalidToken         = errors.New(InvalidToken)
	ErrUnverifiedSubscriber = errors.New(UnverifiedSubscriber)
	ErrSomethingWentWrong   = errors.New(SomethingWentWrong)
	ErrAlreadyVerifiedSub   = errors.New(AlreadyVerifiedSub)
)
