package handlers

import (
	"log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// showMainMenu displays the main menu
func (h *HandlerManager) showMainMenu(userID int64, chatID int64, isMessageSend bool) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Create buttons
	var statusText string
	if userState.IsActive {
		statusText = "Status: 🟢 Online"
	} else {
		statusText = "Status: 🔴 Offline"
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Show Active Users", "show_active"),
			tgbotapi.NewInlineKeyboardButtonData(statusText, "toggle_active"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Settings", "settings"),
			tgbotapi.NewInlineKeyboardButtonData("Find Match", "find_match"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, "Main Menu - Use the buttons below to interact with the bot.")
	msg.ReplyMarkup = keyboard

	if isMessageSend {
		h.bot.Send(msg)
	} else {
		// Edit existing message
		editMsg := tgbotapi.NewEditMessageTextAndMarkup(
			chatID,
			0, // Will be updated in the next step
			"Main Menu - Use the buttons below to interact with the bot.",
			keyboard,
		)
		h.bot.Send(editMsg)
	}
}

// showSettingsMenu displays the settings menu
func (h *HandlerManager) showSettingsMenu(userID int64, chatID int64) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Prepare settings text
	countryText := userState.Settings.Country
	if countryText == "" {
		countryText = "Not set"
	}

	languageText := userState.Settings.Language
	if languageText == "" {
		languageText = "Not set"
	}

	genderText := userState.Settings.Gender
	if genderText == "" {
		genderText = "Not set"
	}

	// Create keyboard
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Country: "+countryText, "set_country"),
			tgbotapi.NewInlineKeyboardButtonData("Clear", "clear_country"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Language: "+languageText, "set_language"),
			tgbotapi.NewInlineKeyboardButtonData("Clear", "clear_language"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Gender: "+genderText, "set_gender"),
			tgbotapi.NewInlineKeyboardButtonData("Clear", "clear_gender"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back to Main Menu", "back_to_main"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		0, // Will be updated by Telegram
		"Settings Menu - Select an option to change or clear:",
		keyboard,
	)

	h.bot.Send(msg)
}

// showLanguageMenu displays language selection menu
func (h *HandlerManager) showLanguageMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("English", "lang_english"),
			tgbotapi.NewInlineKeyboardButtonData("Mandarin", "lang_mandarin"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Hindi", "lang_hindi"),
			tgbotapi.NewInlineKeyboardButtonData("Spanish", "lang_spanish"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("French", "lang_french"),
			tgbotapi.NewInlineKeyboardButtonData("Arabic", "lang_arabic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Bengali", "lang_bengali"),
			tgbotapi.NewInlineKeyboardButtonData("Portuguese", "lang_portuguese"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Russian", "lang_russian"),
			tgbotapi.NewInlineKeyboardButtonData("Japanese", "lang_japanese"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back to Settings", "settings"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		0, // Will be updated by Telegram
		"Select your language:",
		keyboard,
	)

	h.bot.Send(msg)
}

// showGenderMenu displays gender selection menu
func (h *HandlerManager) showGenderMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Male", "gender_male"),
			tgbotapi.NewInlineKeyboardButtonData("Female", "gender_female"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Other", "gender_other"),
			tgbotapi.NewInlineKeyboardButtonData("Back to Settings", "settings"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		0, // Will be updated by Telegram
		"Select your gender:",
		keyboard,
	)

	h.bot.Send(msg)
}

// handleToggleActive toggles a user's active status
func (h *HandlerManager) handleToggleActive(userID int64, chatID int64) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Toggle active status
	userState.IsActive = !userState.IsActive

	// Save updated state
	if err := h.db.SaveUserState(userState); err != nil {
		log.Printf("Error saving user state: %v", err)
		return
	}

	// Show updated menu
	h.showMainMenu(userID, chatID, false)
}

// handleSetSetting sets a user preference
func (h *HandlerManager) handleSetSetting(userID int64, setting string, value string, chatID int64) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Update setting
	switch setting {
	case "country":
		userState.Settings.Country = value
	case "language":
		userState.Settings.Language = value
	case "gender":
		userState.Settings.Gender = value
	}

	// Save updated state
	if err := h.db.SaveUserState(userState); err != nil {
		log.Printf("Error saving user state: %v", err)
		return
	}

	// Show settings menu
	h.showSettingsMenu(userID, chatID)
}

// handleClearSetting clears a user preference
func (h *HandlerManager) handleClearSetting(userID int64, setting string, chatID int64) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Clear setting
	switch setting {
	case "country":
		userState.Settings.Country = ""
	case "language":
		userState.Settings.Language = ""
	case "gender":
		userState.Settings.Gender = ""
	}

	// Save updated state
	if err := h.db.SaveUserState(userState); err != nil {
		log.Printf("Error saving user state: %v", err)
		return
	}

	// Show settings menu
	h.showSettingsMenu(userID, chatID)
}

// checkCompatibility checks if two users are compatible based on their settings
func (h *HandlerManager) checkCompatibility(user1 int64, user2 int64) (bool, error) {
	user1State, err := h.db.GetUserState(user1)
	if err != nil {
		return false, err
	}

	user2State, err := h.db.GetUserState(user2)
	if err != nil {
		return false, err
	}

	// Check gender preference if set
	if user1State.Settings.Gender != "" && user2State.Settings.Gender != "" &&
		user1State.Settings.Gender != user2State.Settings.Gender {
		return false, nil
	}

	// Check language preference if set
	if user1State.Settings.Language != "" && user2State.Settings.Language != "" &&
		user1State.Settings.Language != user2State.Settings.Language {
		return false, nil
	}

	// Check country preference if set
	if user1State.Settings.Country != "" && user2State.Settings.Country != "" &&
		user1State.Settings.Country != user2State.Settings.Country {
		return false, nil
	}

	return true, nil
}

// handleFindMatch tries to find a chat match
func (h *HandlerManager) handleFindMatch(userID int64, chatID int64) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Check if user is active
	if !userState.IsActive {
		h.msgQueue.QueueTextMessage(chatID, "You need to be active to find a match!")
		return
	}

	// Check if user is already in a chat
	if userState.CurrentChat != 0 {
		h.msgQueue.QueueTextMessage(chatID, "You are already in a chat!")
		return
	}

	// Get potential matches
	potentialMatches, err := h.db.FindPotentialMatches(userID)
	if err != nil {
		log.Printf("Error finding potential matches: %v", err)
		h.msgQueue.QueueTextMessage(chatID, "Error finding matches.")
		return
	}

	// Shuffle potential matches
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(potentialMatches), func(i, j int) {
		potentialMatches[i], potentialMatches[j] = potentialMatches[j], potentialMatches[i]
	})

	// Try to find a compatible match
	for _, matchID := range potentialMatches {
		compatible, err := h.checkCompatibility(userID, matchID)
		if err != nil {
			log.Printf("Error checking compatibility: %v", err)
			continue
		}

		if compatible {
			h.startChat(userID, matchID)
			h.msgQueue.QueueTextMessage(chatID, "Match found! Starting chat...")
			return
		}
	}

	// No match found
	h.msgQueue.QueueTextMessage(chatID, "No matches found at the moment. Please try again later.")
}

// startChat starts a chat between two users
func (h *HandlerManager) startChat(user1 int64, user2 int64) error {
	user1State, err := h.db.GetUserState(user1)
	if err != nil {
		return err
	}

	user2State, err := h.db.GetUserState(user2)
	if err != nil {
		return err
	}

	// Update chat states
	user1State.CurrentChat = user2
	user1State.LastActivity = time.Now()
	user2State.CurrentChat = user1
	user2State.LastActivity = time.Now()

	// Save states
	if err := h.db.SaveUserState(user1State); err != nil {
		return err
	}

	if err := h.db.SaveUserState(user2State); err != nil {
		return err
	}

	// Notify users
	h.msgQueue.QueueTextMessage(user1, "Chat started! You can now send messages. Use /end to end the chat.")
	h.msgQueue.QueueTextMessage(user2, "Chat started! You can now send messages. Use /end to end the chat.")

	return nil
}

// handleEndChat ends a chat between two users
func (h *HandlerManager) handleEndChat(userID int64) {
	userState, err := h.db.GetUserState(userID)
	if err != nil {
		log.Printf("Error getting user state: %v", err)
		return
	}

	// Check if user is in a chat
	if userState.CurrentChat == 0 {
		h.msgQueue.QueueTextMessage(userID, "You are not in a chat!")
		return
	}

	partnerID := userState.CurrentChat
	partnerState, err := h.db.GetUserState(partnerID)
	if err != nil {
		log.Printf("Error getting partner state: %v", err)
		return
	}

	// End chat for both users
	userState.CurrentChat = 0
	partnerState.CurrentChat = 0

	// Save states
	if err := h.db.SaveUserState(userState); err != nil {
		log.Printf("Error saving user state: %v", err)
	}

	if err := h.db.SaveUserState(partnerState); err != nil {
		log.Printf("Error saving partner state: %v", err)
	}

	// Notify users
	h.msgQueue.QueueTextMessage(userID, "Chat ended!")
	h.msgQueue.QueueTextMessage(partnerID, "Your chat partner has ended the conversation.")
}
