package TGBOT_CORE

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	TGBOT_ENV "github.com/prokusha/tgbot-golang/internal/telegram-bot/env"
	// "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"

	"log"
	"log/slog"
	"os"
	"time"
)

func TGBotInit() {
	// Get token from the environment variable
	TGBOT_ENV.LoadEnv()

	// Create bot from environment value.
	b, err := gotgbot.NewBot(TGBOT_ENV.Token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
		Logger:      logger,
	})
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{Logger: logger})

	dispatcherHandlers(dispatcher)

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
			// AllowedUpdates: []string{"message", "business_message", "business_connection"},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	logger.Info("Bot has been started...", "bot_username", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
