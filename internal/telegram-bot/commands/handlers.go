package TGBOT_COMMANDS

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	db *pgxpool.Pool
}

func NewHandlers(pool *pgxpool.Pool) *Handlers {
	return &Handlers{db: pool}
}

func (h *Handlers) Ping(b *gotgbot.Bot, ctx *ext.Context) error {
	if err := h.db.Ping(context.Background()); err != nil {
		sendMessage(b, TemplateMessage{
			messageType: TEXT,
			chatId:      ctx.EffectiveSender.Id(),
			data:        "Не прошло",
		})
		return err
	} else {
		sendMessage(b, TemplateMessage{
			messageType: TEXT,
			chatId:      ctx.EffectiveSender.Id(),
			data:        "Прошло))",
		})
	}
	return nil
}
