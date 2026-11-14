package bot

import (
	"context"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"plasitol.dev/iscrav1/internal/messages"
)

// короткое уведомление Showalert
func (shell *Bot) answerCallback(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
	text string,
) {
	if update == nil || update.CallbackQuery == nil {
		shell.logger.Error("[callback] Пустой CallbackQuery")
		return
	}

	query := update.CallbackQuery

	_, err := tgcl.AnswerCallbackQuery(
		ctx,
		&tgbot.AnswerCallbackQueryParams{
			CallbackQueryID: query.ID,
			Text:            text,
		},
	)

	if err != nil {
		shell.logger.Error(
			"[callback] Не удалось ответить на нажатие кнопки",
			"error", err,
		)
	}
}

// текстовый ответ в чат на сообщение
func (shell *Bot) replyCallback(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
	key string,
	data any,
) {
	if update == nil || update.CallbackQuery == nil {
		shell.logger.Error("[callback] Пустой CallbackQuery")
		return
	}

	query := update.CallbackQuery

	text, err := shell.msgs.Get(key, data)
	if err != nil {
		shell.logger.Error(
			"[callback] Не удалось получить текст",
			"key", key,
			"error", err,
		)
		return
	}

	_, err = tgcl.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			ChatID: query.From.ID,
			Text:   text,
		},
	)

	if err != nil {
		shell.logger.Error(
			"[callback] Не удалось отправить сообщение",
			"key", key,
			"error", err,
		)
	}
}

// удаляет сообщение под которым нажали инлайн-кнопку
func (shell *Bot) deleteCallbackMessage(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
) {
	if update == nil || update.CallbackQuery == nil {
		shell.logger.Error("[callback] Пустой CallbackQuery")
		return
	}

	var (
		chatID int64
		msgID  int
	)

	msg := update.CallbackQuery.Message
	switch msg.Type {
	case models.MaybeInaccessibleMessageTypeMessage:
		chatID = msg.Message.Chat.ID
		msgID = msg.Message.ID
	case models.MaybeInaccessibleMessageTypeInaccessibleMessage:
		chatID = msg.InaccessibleMessage.Chat.ID
		msgID = msg.InaccessibleMessage.MessageID
	default:
		shell.logger.Warn("[callback] Нет сообщения для удаления")
		return
	}

	_, err := tgcl.DeleteMessage(ctx, &tgbot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: msgID,
	})
	if err != nil {
		shell.logger.Error(
			"[callback] Не удалось удалить сообщение",
			"error", err,
		)
	}
}

// plasitol
// Дальше ставьте функции, все что системное - выше

// согласие на обработку пд
func (shell *Bot) consentAcceptHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
) {
	if update == nil || update.CallbackQuery == nil {
		shell.logger.Error("[callback] Пустой CallbackQuery")
		return
	}

	shell.logger.Info(
		"[callback] Пользователь нажал кнопку согласия",
	)

	shell.answerCallback( //на ваше усмотрение, надо ли вам это?
		ctx,
		tgcl,
		update,
		messages.Accept, // TODO
	)

	// TODO:
	// сохранить согласие через Consent Service

	shell.deleteCallbackMessage(ctx, tgcl, update)

	shell.replyCallback(
		ctx,
		tgcl,
		update,
		messages.Accept,
		nil,
	)
}

// отказ от обработки персональных данных
func (shell *Bot) consentDeclineHandler(
	ctx context.Context,
	tgcl *tgbot.Bot,
	update *models.Update,
) {
	if update == nil || update.CallbackQuery == nil {
		shell.logger.Error("[callback] Пустой CallbackQuery")
		return
	}

	shell.logger.Info(
		"[callback] Пользователь нажал кнопку отказа",
	)

	shell.answerCallback( //на ваше усмотрение, надо ли вам это?
		ctx,
		tgcl,
		update,
		messages.Decline, // TODO
	)

	shell.deleteCallbackMessage(ctx, tgcl, update)

	// TODO:
	// Сохранить отказ через Consent Service?

	shell.replyCallback(
		ctx,
		tgcl,
		update,
		messages.Decline,
		nil,
	)

}
