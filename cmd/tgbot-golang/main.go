package main

import (
	"log/slog"
	"os"

	. "github.com/prokusha/tgbot-golang/internal/database/core"
	ENV "github.com/prokusha/tgbot-golang/internal/env"
	. "github.com/prokusha/tgbot-golang/internal/telegram-bot/core"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ENV.LoadEnv()

	pool, err := DatabaseInit()
	if err != nil {
		panic("Error init database pool: " + err.Error())
	}
	defer pool.Close()

	tgBot, err := NewTelegramBot(pool)
	if err = tgBot.Start(); err != nil {
		panic("failed to start polling: " + err.Error())
	}

	slog.Info("Bot has been started...", "bot_username", tgBot.GetBotName())
	tgBot.Idle()
}
