package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Constants for the application
const (
	// InactivityTimeout is the duration after which an inactive chat will be terminated
	InactivityTimeout = 1 * time.Hour

	// MatchTimeout is the maximum duration to wait for finding a match
	MatchTimeout = 2 * time.Minute

	// MessageRateLimit is the maximum number of messages per second
	MessageRateLimit = 30
)

// Config holds the application configuration
type Config struct {
	BotToken string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file. Will try to use environment variables directly.")
	}

	// Get bot token from environment variable
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	return &Config{
		BotToken: botToken,
	}
}
