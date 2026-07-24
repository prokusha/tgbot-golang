package TGBOT_CORE

import (
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jackc/pgx/v5/pgxpool"
	ENV "github.com/prokusha/tgbot-golang/internal/env"
	TGBOT_COMMANDS "github.com/prokusha/tgbot-golang/internal/telegram-bot/commands"

	// "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"

	"log"
	"time"
)

type TelegramBot struct {
	b       *gotgbot.Bot
	updater *ext.Updater
}

func (tg *TelegramBot) Start() error {
	return tg.updater.StartPolling(tg.b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
			// AllowedUpdates: []string{"message", "business_message", "business_connection"},
		},
	})
}

func (tg *TelegramBot) Stop() error {
	return tg.updater.Stop()
}

func (tg *TelegramBot) Idle() {
	tg.updater.Idle()
}

func (tg *TelegramBot) GetBotName() string {
	return tg.b.User.Username
}

func NewTelegramBot(pool *pgxpool.Pool) (*TelegramBot, error) {
	// Create bot from environment value.
	b, err := gotgbot.NewBot(ENV.Telegram.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %v", err)
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
		Logger:      slog.Default(),
	})
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{Logger: slog.Default()})

	commHandlers := TGBOT_COMMANDS.NewHandlers(pool)
	dispatcherHandlers(dispatcher, commHandlers)

	return &TelegramBot{b: b, updater: updater}, nil
}
