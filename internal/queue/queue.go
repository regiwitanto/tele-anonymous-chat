package queue

import (
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/regiwitanto/tele-anonymous-chat/internal/config"
	"github.com/regiwitanto/tele-anonymous-chat/internal/models"
)

// MessageQueue manages the queue of messages to be sent
type MessageQueue struct {
	bot     *tgbotapi.BotAPI
	queue   []models.QueuedMessage
	mutex   sync.Mutex
	running bool
	done    chan struct{}
}

// NewMessageQueue creates a new message queue
func NewMessageQueue(bot *tgbotapi.BotAPI) *MessageQueue {
	return &MessageQueue{
		bot:     bot,
		queue:   make([]models.QueuedMessage, 0),
		running: false,
		done:    make(chan struct{}),
	}
}

// Start begins processing the message queue
func (mq *MessageQueue) Start() {
	mq.mutex.Lock()
	if mq.running {
		mq.mutex.Unlock()
		return
	}
	mq.running = true
	mq.mutex.Unlock()

	go mq.processQueue()
}

// Stop stops processing the message queue
func (mq *MessageQueue) Stop() {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if !mq.running {
		return
	}

	mq.running = false
	mq.done <- struct{}{}
}

// QueueTextMessage adds a text message to the queue
func (mq *MessageQueue) QueueTextMessage(chatID int64, text string) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	message := models.QueuedMessage{
		ChatID: chatID,
		Type:   models.TextMessage,
		Text:   text,
	}

	mq.queue = append(mq.queue, message)
}

// QueuePhotoMessage adds a photo message to the queue
func (mq *MessageQueue) QueuePhotoMessage(chatID int64, photoFileID string, caption string) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	message := models.QueuedMessage{
		ChatID:      chatID,
		Type:        models.PhotoMessage,
		PhotoFileID: photoFileID,
		Caption:     caption,
	}

	mq.queue = append(mq.queue, message)
}

// processQueue processes messages from the queue
func (mq *MessageQueue) processQueue() {
	rateLimiter := time.NewTicker(time.Second / config.MessageRateLimit)
	defer rateLimiter.Stop()

	for {
		select {
		case <-mq.done:
			return
		case <-rateLimiter.C:
			// Process one message
			msg, ok := mq.dequeue()
			if !ok {
				continue // Queue is empty
			}

			mq.sendMessage(msg)
		}
	}
}

// dequeue removes and returns the first message in the queue
func (mq *MessageQueue) dequeue() (models.QueuedMessage, bool) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if len(mq.queue) == 0 {
		return models.QueuedMessage{}, false
	}

	msg := mq.queue[0]
	mq.queue = mq.queue[1:]
	return msg, true
}

// sendMessage sends a message based on its type
func (mq *MessageQueue) sendMessage(msg models.QueuedMessage) {
	var err error

	switch msg.Type {
	case models.TextMessage:
		_, err = mq.bot.Send(tgbotapi.NewMessage(msg.ChatID, msg.Text))
	case models.PhotoMessage:
		photoMsg := tgbotapi.NewPhoto(msg.ChatID, tgbotapi.FileID(msg.PhotoFileID))
		if msg.Caption != "" {
			photoMsg.Caption = msg.Caption
		}
		_, err = mq.bot.Send(photoMsg)
	}

	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
