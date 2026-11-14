package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// единственное место где текст из меседжа идет в телегу
func (shell *Bot) sendKey(
	ctx context.Context,
	tgcl *tgbot.Bot,
	chatID int64,
	key string,
	data any,
	markup models.ReplyMarkup,
) {
	text, err := shell.msgs.Get(key, data)
	if err != nil {
		shell.logger.Error("Нет текста", "key", key, "error", err)
		return
	}

	_, err = tgcl.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: markup,
	})
	if err != nil {
		shell.logger.Error("Ошибка отправки",
			"key", key, "chat_id", chatID, "error", err)
	}
}

// ответ на обычное сообщение
func (shell *Bot) replyHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
	key string,
	data any,
	markup models.ReplyMarkup,
) {
	if update.Message == nil {
		return
	}
	shell.sendKey(ctx, tgcl, update.Message.Chat.ID, key, data, markup)
}
