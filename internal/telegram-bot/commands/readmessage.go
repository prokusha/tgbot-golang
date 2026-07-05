package TGBOT_COMMANDS

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	TEXT_HELPER "github.com/prokusha/tgbot-golang/internal/helper/text"
	MUSIC_SERVICES "github.com/prokusha/tgbot-golang/internal/music-services"
)

func analyzeMessage(msg *gotgbot.Message) (bool, TemplateMessage) {
	analyze := TEXT_HELPER.AnalyzeTextType(msg.GetText())
	var message TemplateMessage
	check := true
	switch analyze {
	case TEXT_HELPER.TextEvent:
		message.messageType = TEXT
		data := fmt.Sprintf("Пользователь @%s оставил сообщение с датой.\n<blockquote>%s</blockquote>\nСоздать напоминание?", msg.GetSender().Username(), msg.GetText())
		message.data = data
		message.buttons = YesNoButtons
	case TEXT_HELPER.TextURL:
		if url, typeUrl := TEXT_HELPER.URLAndType(msg.GetText()); typeUrl != TEXT_HELPER.URLNull {
			message.messageType = AUDIO
			data, _ := MUSIC_SERVICES.GetMusicURL(url[0])
			message.data = data
		}
	default:
		check = false
	}
	return check, message
}

func ReadMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	if msg := ctx.Update.BusinessMessage; msg != nil {
		if ok, message := analyzeMessage(msg); ok {
			message.chatId = msg.Chat.Id
			message.assistentId = msg.BusinessConnectionId
			return sendMessage(b, message)
		}
	} else {
		if ok, message := analyzeMessage(ctx.Message); ok {
			message.chatId = ctx.Message.Chat.Id
			return sendMessage(b, message)
		}
	}
	return nil
}
