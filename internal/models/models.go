package models

import (
	"time"
)

// UserState represents the state of a user in the system
type UserState struct {
	UserID         int64
	IsActive       bool
	CurrentChat    int64
	LastActivity   time.Time
	Settings       UserSettings
	MatchStartTime *time.Time
}

// UserSettings contains user preferences for matching
type UserSettings struct {
	Country  string
	Language string
	Gender   string
}

// NewUserState creates a new UserState instance
func NewUserState(userID int64) *UserState {
	return &UserState{
		UserID:       userID,
		IsActive:     false,
		CurrentChat:  0,
		LastActivity: time.Now(),
		Settings: UserSettings{
			Country:  "",
			Language: "",
			Gender:   "",
		},
		MatchStartTime: nil,
	}
}

// ToMap converts a UserState to a map for database storage
func (u *UserState) ToMap() map[string]interface{} {
	lastActivity := u.LastActivity.Format(time.RFC3339)

	return map[string]interface{}{
		"is_active":     u.IsActive,
		"current_chat":  u.CurrentChat,
		"last_activity": lastActivity,
		"country":       u.Settings.Country,
		"language":      u.Settings.Language,
		"gender":        u.Settings.Gender,
	}
}

// MessageType represents the type of message to be sent
type MessageType int

const (
	// TextMessage is a simple text message
	TextMessage MessageType = iota

	// PhotoMessage is a photo with optional caption
	PhotoMessage
)

// QueuedMessage represents a message in the queue to be sent
type QueuedMessage struct {
	ChatID      int64
	Type        MessageType
	Text        string
	PhotoFileID string
	Caption     string
}
