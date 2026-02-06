# Go Bot Setup Guide

Place your Go files in this directory (`go-bot/`).

## How to run:
1. Ensure Go is installed on your VPS.
2. Build the bot: `go build -o discord-bot main.go`
3. Run with PM2: `pm2 start ./discord-bot --name "discord-bot"`

The main application will handle the frontend and API, while this bot can run as a background process managed by PM2.
