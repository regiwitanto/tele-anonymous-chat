package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/regiwitanto/tele-anonymous-chat/internal/bot"
	"github.com/regiwitanto/tele-anonymous-chat/internal/config"
	"github.com/regiwitanto/tele-anonymous-chat/internal/database"
)

func main() {
	// Initialize logging
	log.SetOutput(os.Stdout)
	log.Println("Starting Telegram Anonymous P2P Chat Bot...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDB("database.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized")

	// Create bot instance
	telegramBot, err := bot.NewBot(cfg, db)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	log.Println("Bot created successfully")

	// Start bot in a separate goroutine
	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatalf("Error starting bot: %v", err)
		}
	}()
	log.Println("Bot started successfully")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down bot...")
	telegramBot.Stop()
	log.Println("Bot stopped successfully")
}
