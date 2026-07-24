package TGBOT_CORE

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	. "github.com/prokusha/tgbot-golang/internal/telegram-bot/commands"
)

func dispatcherHandlers(dispatcher *ext.Dispatcher, commHandlers *Handlers) {
	dispatcher.AddHandler(handlers.NewCommand("ping", commHandlers.Ping))
	dispatcher.AddHandler(handlers.NewCommand("start", Start))
	// dispatcher.AddHandler(handlers.NewCommand("get_music", GetMusic))
	dispatcher.AddHandler(handlers.NewCommand("search", Search))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("get_music:"), GetMusic))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("search_result:"), CB_search_result))
	dispatcher.AddHandler(handlers.BusinessConnection{
		Filter: func(bc *gotgbot.BusinessConnection) bool {
			return bc.IsEnabled
		},
		Response: BusinessConnectionHandler,
	})
	dispatcher.AddHandler(handlers.NewMessage(nil, ReadMessage).SetAllowBusiness(true))
	// dispatcher.AddHandler(handlers.NewMessage(nil, DaPizda).SetAllowBusiness(true))
}
