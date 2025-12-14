# University Confession Bot

## Overview
A Telegram bot for university confessions, featuring:
- Anonymous confession submissions
- Admin moderation in a private group
- Blind chat functionality (random user pairing)
- Voting system
- Scheduled posting
- Content moderation with banned word filtering

## Tech Stack
- **Language**: Go 1.24
- **Database**: SQLite (file-based, `confessions.db`)
- **Telegram API**: go-telegram-bot-api/telegram-bot-api v5
- **Mode**: Webhook-based (listens on port 8080)

## Project Structure
- `main.go` - Entry point, bot initialization, webhook server
- `config.go` - Configuration constants (admin IDs, channel IDs, timing)
- `database.go` - SQLite database initialization and operations
- `handlers.go` - Telegram message and callback handlers
- `admin.go` - Admin command handlers
- `moderation.go` - Content moderation logic
- `voting.go` - Voting system for confessions
- `blindchat.go` - Anonymous chat pairing feature
- `schedular.go` - Scheduled task management
- `webhook.go` - Webhook endpoint handler

## Environment Variables
- `BOT_TOKEN` - Telegram Bot API token (required)

## Configuration
Edit `config.go` to set:
- `AdminID` - Telegram user ID of the bot owner
- `AdminGroupID` - Group ID for admin moderation
- `ChannelID` - Channel ID for posting confessions

## Running
```bash
go run .
```

## Deployment
The bot runs as a webhook server on port 8080. Configure your domain's webhook URL in Telegram using the BotFather or API.
