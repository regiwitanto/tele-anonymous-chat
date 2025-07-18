package handlers

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/regiwitanto/tele-anonymous-chat/internal/database"
	"github.com/regiwitanto/tele-anonymous-chat/internal/queue"
)

// HandlerManager manages all the telegram update handlers
type HandlerManager struct {
	bot      *tgbotapi.BotAPI
	db       *database.DB
	msgQueue *queue.MessageQueue
}

// NewHandlerManager creates a new handler manager
func NewHandlerManager(bot *tgbotapi.BotAPI, db *database.DB, msgQueue *queue.MessageQueue) *HandlerManager {
	return &HandlerManager{
		bot:      bot,
		db:       db,
		msgQueue: msgQueue,
	}
}

// HandleCommand processes command messages
func (h *HandlerManager) HandleCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	userID := update.Message.From.ID

	switch command {
	case "start":
		h.handleStart(update)
	case "end":
		h.handleEndChat(userID)
	default:
		h.msgQueue.QueueTextMessage(update.Message.Chat.ID, "Unknown command. Use /start to see available options.")
	}
}

// HandleCallback processes callback queries (button clicks)
func (h *HandlerManager) HandleCallback(update tgbotapi.Update) {
	query := update.CallbackQuery
	userID := query.From.ID
	callbackData := query.Data

	// Send an empty callback response to stop the loading animation
	callbackConfig := tgbotapi.NewCallback(query.ID, "")
	h.bot.Send(callbackConfig)

	switch callbackData {
	case "show_active":
		h.handleShowActive(query.Message.Chat.ID)

	case "toggle_active":
		h.handleToggleActive(userID, query.Message.Chat.ID)

	case "settings":
		h.showSettingsMenu(userID, query.Message.Chat.ID)

	case "back_to_main":
		h.showMainMenu(userID, query.Message.Chat.ID, false)

	case "find_match":
		h.handleFindMatch(userID, query.Message.Chat.ID)

	case "set_country":
		h.msgQueue.QueueTextMessage(query.Message.Chat.ID, "Please enter your country (e.g., USA, UK, etc.):")
		// We'd need to handle text input in a stateful way

	case "clear_country":
		h.handleClearSetting(userID, "country", query.Message.Chat.ID)

	case "set_language":
		h.showLanguageMenu(query.Message.Chat.ID)

	case "clear_language":
		h.handleClearSetting(userID, "language", query.Message.Chat.ID)

	case "set_gender":
		h.showGenderMenu(query.Message.Chat.ID)

	case "clear_gender":
		h.handleClearSetting(userID, "gender", query.Message.Chat.ID)
	}

	// Handle language selection
	if len(callbackData) > 5 && callbackData[:5] == "lang_" {
		language := callbackData[5:]
		h.handleSetSetting(userID, "language", language, query.Message.Chat.ID)
	}

	// Handle gender selection
	if len(callbackData) > 7 && callbackData[:7] == "gender_" {
		gender := callbackData[7:]
		h.handleSetSetting(userID, "gender", gender, query.Message.Chat.ID)
	}
}

// HandleMessage processes regular messages
func (h *HandlerManager) HandleMessage(update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Get user state
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// If user is in a chat, forward the message to their chat partner
	if userState.CurrentChat > 0 {
		// Update last activity
		userState.LastActivity = time.Now()
		if err := h.db.SaveUserState(userState); err != nil {
			log.Printf("Error saving user state: %v", err)
		}

		if update.Message.Photo != nil {
			// Handle photo messages
			photos := update.Message.Photo
			// Get the largest available photo
			photoFileID := photos[len(photos)-1].FileID
			caption := "Anonymous sent a photo"
			if update.Message.Caption != "" {
				caption = "Anonymous: " + update.Message.Caption
			}
			h.msgQueue.QueuePhotoMessage(userState.CurrentChat, photoFileID, caption)
		} else {
			// Handle text messages
			h.msgQueue.QueueTextMessage(userState.CurrentChat, fmt.Sprintf("Anonymous: %s", update.Message.Text))
		}
	} else {
		// If not in a chat, show main menu
		h.showMainMenu(userID, chatID, false)
	}
}

// EndInactiveChats terminates chats that have been inactive for too long
func (h *HandlerManager) EndInactiveChats() error {
	// Implementation will be in a separate file
	return h.checkAndEndInactiveChats()
}

// handleStart handles the /start command
func (h *HandlerManager) handleStart(update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Welcome message
	welcomeMsg := `Welcome to the Anonymous P2P Chat Bot!

How it works:
- This bot lets you chat anonymously with random users.
- You can set preferences (country, language, gender) to match with similar users.
- Only text and photo messages are allowed.
- Chats are ended automatically after 1 hour of inactivity.

Commands and Features:
/start - Show this message and the main menu.
/end - End your current anonymous chat.
Show Active Users - See how many users are currently online.
Status: Online/Offline - Toggle your availability for matching.
Settings - Set or clear your country, language, or gender preferences.
Find Match - Start searching for a random chat partner.

Use the menu buttons to navigate. Enjoy chatting!`

	h.msgQueue.QueueTextMessage(chatID, welcomeMsg)

	// Show main menu
	h.showMainMenu(userID, chatID, true)
}

// handleShowActive shows the number of active users
func (h *HandlerManager) handleShowActive(chatID int64) {
	count, err := h.db.GetActiveUsers()
	if err != nil {
		h.msgQueue.QueueTextMessage(chatID, "Error getting active users count.")
		return
	}

	h.msgQueue.QueueTextMessage(chatID, fmt.Sprintf("Active users: %d", count))
}
