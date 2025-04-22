# Telegram Anonymous P2P Chat Bot

A privacy-focused Telegram bot that enables anonymous peer-to-peer chatting with smart matching based on user preferences. Connect with random users while maintaining your privacy and finding matches based on your preferences.

## Features

- 🔒 **Anonymous Chatting**: Chat with random users while keeping your identity private
- 🎯 **Smart Matching**: Get matched with users based on your preferences:
  - Country
  - Language
  - Gender
- ⚡ **Real-time Status**: See who's online and available to chat
- 🖼️ **Media Support**: Send and receive photos in chats
- ⏱️ **Auto Timeouts**: 
  - Chat inactivity timeout (1 hour)
  - Match search timeout (2 minutes)
- 🔄 **Rate Limiting**: Respects Telegram's API rate limits
- ⚙️ **Customizable Settings**: Set and clear your preferences anytime

## Installation

### Prerequisites

- Python 3.7 or higher
- A Telegram Bot Token (get it from [@BotFather](https://t.me/BotFather))
- SQLite3 (usually comes with Python)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/okoyausman/telegram-anonymous-p2p-chat.git
cd telegram-anonymous-p2p-chat
```

2. Install required packages:
```bash
pip install python-telegram-bot python-dotenv
```

3. Create a `.env` file in the project root:
```bash
BOT_TOKEN=your_telegram_bot_token_here
```

4. Run the bot:
```bash
python main.py
```

## Usage

1. Start the bot by sending `/start` command
2. Use the main menu to:
   - Toggle your online status
   - View active users
   - Access settings
   - Find a match

3. In Settings, you can:
   - Set your country
   - Choose your language
   - Set your gender
   - Clear any preference

4. When in a chat:
   - Send text messages and photos
   - Use `/end` to end the chat
   - Chat will auto-end after 1 hour of inactivity

## Configuration

The bot has several configurable constants in `main.py`:

```python
INACTIVITY_TIMEOUT = 3600  # 1 hour in seconds
MATCH_TIMEOUT = 120  # 2 minutes in seconds
MESSAGE_RATE_LIMIT = 30  # messages per second
```

## Database

The bot uses SQLite3 to store:
- User states
- Chat connections
- User preferences
- Activity timestamps

## Security & Privacy

- All chats are anonymous
- No user data is stored beyond necessary preferences
- Messages are not logged
- Users can clear their preferences anytime

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

If you encounter any issues or have questions, please open an issue in the GitHub repository. 
