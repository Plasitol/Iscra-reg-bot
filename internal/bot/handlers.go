package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"plasitol.dev/iscrav1/internal/messages"
)

func (shell *Bot) replyHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
	key string,
	data any,
	keyboard models.ReplyMarkup,
) {
	user := update.Message.From

	text, err := shell.msgs.Get(key, data)
	if err != nil {
		shell.logger.Error(
			"[handler] Не удалось получить текст",
			"key", key,
			"error", err)
		return
	}
	_, err = tgcl.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		shell.logger.Error(
			"Ошибка отправки", "key", key,
			"error", err,
			"telegram_id", user.ID,
			"username", user.Username,
		)
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

	shell.replyHandler(ctx, tgcl, update, messages.Start, map[string]any{
		"Name": user.FirstName,
	}, mainKeyboard())
}

// обработка команды (Регистрация)
func (shell *Bot) registrationHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	shell.replyHandler(ctx, tgcl, update, messages.Registration, nil, consentKeyboard())
}

// обработка команды (Помощь)
func (shell *Bot) helpHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	shell.replyHandler(ctx, tgcl, update, messages.Help, nil, nil)
}

// обработка дефолта
func (shell *Bot) defaultHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	shell.replyHandler(ctx, tgcl, update, messages.Unknown, nil, nil) // можете просто заменить на Help
}
