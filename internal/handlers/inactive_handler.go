package handlers

import (
	"log"
	"time"

	"github.com/regiwitanto/tele-anonymous-chat/internal/config"
)

// checkAndEndInactiveChats terminates chats that have been inactive for too long
func (h *HandlerManager) checkAndEndInactiveChats() error {
	// Get all active chats
	activeChats, err := h.db.GetActiveChats()
	if err != nil {
		return err
	}

	now := time.Now()

	for _, chat := range activeChats {
		// Find the most recent activity
		lastActivity := chat.LastActivity1
		if chat.LastActivity2.After(chat.LastActivity1) {
			lastActivity = chat.LastActivity2
		}

		// Check if chat is inactive
		if now.Sub(lastActivity) > config.InactivityTimeout {
			// End the chat due to inactivity
			user1State, err := h.db.GetUserState(chat.User1ID)
			if err != nil {
				log.Printf("Error getting user state: %v", err)
				continue
			}

			user2State, err := h.db.GetUserState(chat.User2ID)
			if err != nil {
				log.Printf("Error getting user state: %v", err)
				continue
			}

			// Clear chat states
			user1State.CurrentChat = 0
			user2State.CurrentChat = 0

			// Save updated states
			if err := h.db.SaveUserState(user1State); err != nil {
				log.Printf("Error saving user state: %v", err)
			}

			if err := h.db.SaveUserState(user2State); err != nil {
				log.Printf("Error saving user state: %v", err)
			}

			// Notify users
			h.msgQueue.QueueTextMessage(chat.User1ID, "Chat ended due to inactivity!")
			h.msgQueue.QueueTextMessage(chat.User2ID, "Chat ended due to inactivity!")
		}
	}

	return nil
}
