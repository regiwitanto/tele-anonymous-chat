package utils

import (
	"time"

	"github.com/regiwitanto/tele-anonymous-chat/internal/models"
)

// CheckInactiveTimeout checks if a user's last activity exceeds the timeout
func CheckInactiveTimeout(userState *models.UserState, timeout time.Duration) bool {
	return time.Since(userState.LastActivity) > timeout
}

// CheckMatchTimeout checks if a user has been waiting for a match too long
func CheckMatchTimeout(userState *models.UserState, timeout time.Duration) bool {
	if userState.MatchStartTime == nil {
		return false
	}
	return time.Since(*userState.MatchStartTime) > timeout
}

// FormatTimestamp formats a time.Time into a readable string
func FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTimestamp parses a string timestamp into time.Time
func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}
