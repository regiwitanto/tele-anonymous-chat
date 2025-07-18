package bot

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/regiwitanto/tele-anonymous-chat/internal/config"
	"github.com/regiwitanto/tele-anonymous-chat/internal/database"
	"github.com/regiwitanto/tele-anonymous-chat/internal/handlers"
	"github.com/regiwitanto/tele-anonymous-chat/internal/queue"
)

// Bot represents the Telegram bot
type Bot struct {
	api      *tgbotapi.BotAPI
	db       *database.DB
	msgQueue *queue.MessageQueue
	handlers *handlers.HandlerManager
	config   *config.Config
	stopChan chan struct{}
}

// NewBot creates a new Bot instance
func NewBot(cfg *config.Config, db *database.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	msgQueue := queue.NewMessageQueue(api)

	bot := &Bot{
		api:      api,
		db:       db,
		msgQueue: msgQueue,
		config:   cfg,
		stopChan: make(chan struct{}),
	}

	bot.handlers = handlers.NewHandlerManager(api, db, msgQueue)

	return bot, nil
}

// Start starts the bot and begins processing updates
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	// Start message queue processing
	b.msgQueue.Start()

	// Start inactivity checker
	go b.checkInactiveChats()

	// Configure update channel
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.api.GetUpdatesChan(updateConfig)

	// Process updates
	for {
		select {
		case update := <-updates:
			go b.handleUpdate(update)
		case <-b.stopChan:
			return nil
		}
	}
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.msgQueue.Stop()
	b.stopChan <- struct{}{}
}

// handleUpdate processes an incoming update
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	// Handle commands
	if update.Message != nil && update.Message.IsCommand() {
		b.handlers.HandleCommand(update)
		return
	}

	// Handle callback queries (button clicks)
	if update.CallbackQuery != nil {
		b.handlers.HandleCallback(update)
		return
	}

	// Handle messages
	if update.Message != nil {
		b.handlers.HandleMessage(update)
		return
	}
}

// checkInactiveChats periodically checks for inactive chats
func (b *Bot) checkInactiveChats() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := b.handlers.EndInactiveChats(); err != nil {
				log.Printf("Error ending inactive chats: %v", err)
			}
		case <-b.stopChan:
			return
		}
	}
}
