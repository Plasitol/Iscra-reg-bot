package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

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

	_, err := tgcl.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "Привет!\n\n" +
			"барамбимбимбум", // TODO! ТЕКСТ
	})

	if err != nil {
		shell.logger.Error(
			"Ошибка отправки сообщения",
			"error", err,
			"telegram_id", user.ID,
		)
	}
}

// обработка команды /help
func (shell *Bot) helpHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
) {
	_, err := tgcl.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "Доступные команды:\n\n" + // TODO! ТЕКСТ
			"/start начать работу\n" +
			"/help помощь",
	})

	if err != nil {
		shell.logger.Error("Ошибка отправки help", "error", err)
	}
}

// обработка дефолта
func (shell *Bot) defaultHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
) {
	if update.Message == nil {
		return
	}

	_, err := tgcl.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Ниче не понимаю", // TODO! ТЕКСТ
	})

	if err != nil {
		shell.logger.Error("Ошибка defaultHandler", "error", err)
	}
}
