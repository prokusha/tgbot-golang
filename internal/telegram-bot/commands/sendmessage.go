package TGBOT_COMMANDS

import (
	"fmt"
	"log"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type MessageType int

const (
	AUDIO MessageType = iota
	TEXT
)

type TemplateMessage struct {
	messageType MessageType
	chatId      int64
	assistentId string
	data        string
	buttons     [][]gotgbot.InlineKeyboardButton
}

var YesNoButtons = [][]gotgbot.InlineKeyboardButton{
	{
		{
			Text:         "Да",
			CallbackData: "add_event:",
		},
		{
			Text:         "Нет",
			CallbackData: "delete_message:",
		},
	},
}

func sendMessage(b *gotgbot.Bot, message TemplateMessage) error {
	var err error
	switch message.messageType {
	case TEXT:
		err = sendTextMessage(b, message)
	case AUDIO:
		err = sendAudioMessage(b, message)
	}
	return err
}

func sendTextMessage(b *gotgbot.Bot, message TemplateMessage) error {
	_, err := b.SendMessage(message.chatId, message.data, &gotgbot.SendMessageOpts{
		BusinessConnectionId: message.assistentId,
		ParseMode:            "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: message.buttons,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

func sendAudioMessage(b *gotgbot.Bot, message TemplateMessage) error {
	log.Printf("Try to send file by path: %s\n", message.data)
	file, err := os.Open(message.data)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = b.SendAudio(message.chatId, gotgbot.InputFileByReader("", file), &gotgbot.SendAudioOpts{
		BusinessConnectionId: message.assistentId,
		ParseMode:            "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: message.buttons,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	// os.Remove(message.data)
	return nil
}
