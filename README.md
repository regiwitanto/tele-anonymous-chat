# Telegram Anonymous Chat Bot

A privacy-focused Telegram bot that enables anonymous peer-to-peer chatting with smart matching based on user preferences.

## Features

- 🔒 **Anonymous Chatting**: Chat with random users while keeping your identity private
- 🎯 **Smart Matching**: Match with users based on country, language, and gender preferences
- ⚡ **Real-time Status**: See who's online and available to chat
- 🖼️ **Media Support**: Send and receive photos in chats
- ⏱️ **Auto Timeouts**: Inactive chats end after 1 hour, matching timeout after 2 minutes
- 🔄 **Rate Limiting**: Respects Telegram API limits
- ⚙️ **Customizable Settings**: Set and clear your preferences anytime

## Quick Start

### Prerequisites

- Go 1.20 or higher
- A Telegram Bot Token from [@BotFather](https://t.me/BotFather)
- SQLite3

### Setup

1. Clone the repository:
```bash
git clone https://github.com/regiwitanto/tele-anonymous-chat.git
cd tele-anonymous-chat
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create a `.env` file:
```bash
BOT_TOKEN=your_telegram_bot_token_here
```

4. Run the bot:
```bash
go run main.go
```

5. Or build and run:
```bash
go build
./tele-anonymous-chat
```

## Usage

1. Start the bot with `/start` command
2. Main Menu Options:
   - Toggle your online status
   - View active users
   - Access settings
   - Find a match
3. Send text and photos in chats
4. Use `/end` to end conversations

## Project Structure

```
tele-anonymous-chat/
├── cmd/bot/          # Application entry point
├── internal/         # Internal packages
│   ├── bot/          # Bot functionality
│   ├── config/       # App configuration
│   ├── database/     # Database operations
│   ├── handlers/     # Message handlers
│   ├── models/       # Data models
│   ├── queue/        # Message queue
│   └── utils/        # Utilities
├── main.go           # Main entry point
├── go.mod            # Go module definition
├── .env              # Environment variables
└── README.md         # Documentation
```

## Configuration

Key settings in `internal/config/config.go`:

```go
// Chat timeouts and rate limits
InactivityTimeout = 1 * time.Hour
MatchTimeout = 2 * time.Minute
MessageRateLimit = 30
```

## Database

SQLite3 stores:
- User states and preferences
- Chat connections
- Activity timestamps

## Features

### Privacy
- All chats are anonymous
- Only necessary preferences are stored
- Messages are not logged

### Technical
- Written in Go for performance
- Concurrent message handling
- Rate-limited message queue
- Automatic timeout handling

## License

MIT License

## Support

For issues or questions, please open an issue on GitHub. 
