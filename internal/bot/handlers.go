package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"plasitol.dev/iscrav1/internal/messages"
)

func (shell *Bot) reply(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
	key string,
	data any,
) {
	user := update.Message.From

	text, err := shell.msgs.Get(key, data)
	if err != nil {
		shell.logger.Error("Не удалось получить текст", "key", key, "error", err)
		return
	}
	_, err = tgcl.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		shell.logger.Error(
			"Ошибка отправки", "key", key,
			"error", err,
			"telegram_id", user.ID,
			"username", user.Username)
	}
}

// обработка команды /start
func (shell *Bot) startHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
) {
	user := update.Message.From

	shell.logger.Info(
		"Юзер запустил бота",
		"telegram_id", user.ID,
		"username", user.Username,
	)

	shell.reply(ctx, tgcl, update, messages.Start, map[string]any{
		"Name": user.FirstName,
	})
}

// обработка команды /help
func (shell *Bot) helpHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	shell.reply(ctx, tgcl, update, messages.Help, nil)
}

// обработка дефолта
func (shell *Bot) defaultHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	shell.reply(ctx, tgcl, update, messages.Unknown, nil) // можете просто заменить на Help
}
