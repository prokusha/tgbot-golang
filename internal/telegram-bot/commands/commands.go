package TGBOT_COMMANDS

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	. "github.com/prokusha/tgbot-golang/internal/music-services"

	"fmt"
	"log"
	"strconv"
	"strings"
)

var ResultCache = make(map[int64][]MusicData)

func BusinessConnectionHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	conn := ctx.Update.BusinessConnection
	log.Println(conn.User.Id)
	return nil
}

func daToPizda(original_text string) (bool, string, int, int) {
	m := strings.ToLower(original_text)
	m = strings.ReplaceAll(m, ",", " ")
	m = strings.TrimSpace(m)
	runes := []rune(m)
	for i := 0; i < len(runes)-1; i++ {
		if (runes[i] == 'd' || runes[i] == 'д') && (runes[i+1] == 'a' || runes[i+1] == 'а') {
			map_pizda := map[rune]string{
				'd': "Piz",
				'д': "Пиз",
				'a': "da",
				'а': "да",
			}
			var answer strings.Builder
			answer.WriteString(map_pizda[runes[i]])
			j := i + 1
			for j < len(runes) && (runes[j] == 'a' || runes[j] == 'а') {
				answer.WriteString(map_pizda[runes[j]])
				j++
			}
			answer.WriteString("))))")
			return true, answer.String(), i, j
		}
	}
	return false, "", 0, 0
}

func DaPizda(b *gotgbot.Bot, ctx *ext.Context) error {
	// log.Println("Message", *ctx.Update)
	if ctx.Update.BusinessMessage != nil {
		msg := ctx.Update.BusinessMessage
		connId := msg.BusinessConnectionId
		check, answer, x, y := daToPizda(msg.Text)
		if check == false {
			return nil
		}
		log.Printf("Бизнес-текст: %s от %d\n", msg.Text, msg.From.Id)
		_, err := b.SendMessage(msg.Chat.Id, answer, &gotgbot.SendMessageOpts{
			BusinessConnectionId: connId,
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: msg.MessageId,
				Quote:     string([]rune(msg.Text)[x:y]),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to echo message: %w", err)
		}
	} else {
		original_text := ctx.Message.GetText()
		check, answer, x, y := daToPizda(original_text)
		if check == false {
			return nil
		}
		log.Printf("Просто текст: %s\n", original_text)
		_, err := b.SendMessage(ctx.Message.Chat.Id, answer, &gotgbot.SendMessageOpts{
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: ctx.Message.MessageId,
				Quote:     string([]rune(original_text)[x:y]),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to echo message: %w", err)
		}
	}
	return nil
}

func GetMusic(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println("GET_MUSIC")
	cb := ctx.Update.CallbackQuery

	// Обязательно подтверждаем callback, чтобы убрать индикатор загрузки у пользователя
	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Загрузка...",
	})
	if err != nil {
		return err
	}

	data := cb.Data
	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		return nil
	}

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to conver query to int: %w", err)
	}

	// response, err := http.Get("https:" + ResultCache[ctx.EffectiveChat.Id][index].URL)
	// if err != nil {
	// 	return fmt.Errorf("failed to get mp3: %w", err)
	// }
	// defer response.Body.Close()
	log.Printf("Выбранна композиция: %s-%s", ResultCache[ctx.EffectiveChat.Id][index].Author, ResultCache[ctx.EffectiveChat.Id][index].Name)
	_, err = b.SendDocument(ctx.EffectiveChat.Id, gotgbot.InputFileByURL("https:"+ResultCache[ctx.EffectiveChat.Id][index].URL), &gotgbot.SendDocumentOpts{
		// Caption:   "Вот ваш запрашиваемый файл! 📁",
		ParseMode: "HTML",
	})

	if err != nil {
		return fmt.Errorf("failed to send source: %w", err)
	}
	return nil
}

func result_list(index int, id int64) [][]gotgbot.InlineKeyboardButton {
	inline := [][]gotgbot.InlineKeyboardButton{}
	for i := index; i < index+10 && i < len(ResultCache[id]); i++ {
		list := ResultCache[id]
		inline = append(inline, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%s - %s", list[i].Author, list[i].Name),
				CallbackData: fmt.Sprintf("get_music:%d", i),
			},
		})
	}
	arrows := []gotgbot.InlineKeyboardButton{}
	arrows = append(arrows, gotgbot.InlineKeyboardButton{
		Text:         "⬅️",
		CallbackData: fmt.Sprintf("search_result:%d", index-10),
	})
	arrows = append(arrows, gotgbot.InlineKeyboardButton{
		Text:         "➡️",
		CallbackData: fmt.Sprintf("search_result:%d", index+10),
	})

	inline = append(inline, arrows)

	return inline
}

func CB_search_result(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Загрузка...",
	})
	if err != nil {
		return err
	}

	data := cb.Data
	parts := strings.Split(data, ":")
	index := 0
	if len(parts) == 2 {
		index, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("failed to conver query to int: %w", err)
		}
		if index < 0 || index > len(ResultCache[ctx.EffectiveChat.Id]) {
			return nil
		}
	}

	_, _, err = cb.Message.EditText(b, fmt.Sprintf("Результат поиска: %d-%d", index+1, min(index+10, len(ResultCache[ctx.EffectiveChat.Id]))), &gotgbot.EditMessageTextOpts{
		ParseMode:   "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: result_list(index, ctx.EffectiveChat.Id)},
	})
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

func Search(b *gotgbot.Bot, ctx *ext.Context) error {
	ResultCache[ctx.EffectiveChat.Id] = nil
	text := ctx.Message.GetText()
	query := strings.Trim(text, "/search")
	query = strings.TrimSpace(query)
	log.Println(query)
	if query == "" {
		_, err := ctx.EffectiveMessage.Reply(b, "Пустая строка", nil)
		if err != nil {
			return fmt.Errorf("failed to echo message: %w", err)
		}
		return nil
	}
	music_list, err := SearchMusic(query)
	if err != nil {
		return fmt.Errorf("failed to search music: %w", err)
	}
	ResultCache[ctx.EffectiveChat.Id] = music_list
	_, err = ctx.Message.Reply(b, "Результат поиска: ", &gotgbot.SendMessageOpts{
		ParseMode:   "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: result_list(0, ctx.EffectiveChat.Id)},
	})
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}

// start introduces the bot.
func Start(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Println(ctx.Message.GetText())
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s.\nI am a sample bot to demonstrate how file sending works.\n\nTry the /source command!", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Search", CallbackData: "search"},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func CBLOG(b *gotgbot.Bot, ctx *ext.Context) error {
	cb := ctx.Update.CallbackQuery

	_, err := cb.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: "Загрузка...",
	})
	if err != nil {
		return err
	}

	data := cb.Data
	log.Println(data)
	return nil
}
