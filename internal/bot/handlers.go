package bot

import (
	"context"
	"errors"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"plasitol.dev/iscrav1/internal/messages"
	registration "plasitol.dev/iscrav1/internal/registation" // TODO регистейшон, лол, поправить надо
	"plasitol.dev/iscrav1/internal/user"
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

	shell.replyHandler(ctx, tgcl, update, messages.Start, map[string]any{
		"Name": user.FirstName,
	}, mainKeyboard())
}

// обработка команды (Регистрация)
func (shell *Bot) registrationHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	from := update.Message.From

	done, err := shell.users.IsRegistered(ctx, from.ID)
	if err != nil {
		shell.logger.Error(
			"Проверка регистрации",
			"error", err,
			"telegram_id", from.ID,
			"username", from.Username)
		shell.replyHandler(ctx, tgcl, update, messages.TryLater, nil, nil)
		return
	}
	if done {
		shell.replyHandler(ctx, tgcl, update, messages.AlreadyRegistered, nil, nil)
		return
	}

	shell.replyHandler(ctx, tgcl, update, messages.Registration, nil, consentKeyboard())
}

// обработка согласия на регистрацию (Регистрация)
func (shell *Bot) confirmHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}
	from := update.CallbackQuery.From

	shell.answerCallback(ctx, tgcl, update, "")

	sess, ok := shell.sessions.Get(from.ID)
	if !ok {
		shell.deleteCallbackMessage(ctx, tgcl, update)
		shell.sendKey(ctx, tgcl, from.ID, messages.Unknown, nil, nil)
		return // сессии нет (рестарт бота или двойное нажатие) TODO
	}

	a := sess.Answers
	_, err := shell.users.Register(ctx, user.User{
		TelegramID:   from.ID,
		Username:     from.Username,
		FullName:     a["fio"],
		Passport:     a["passport"],
		BirthDate:    a["birth_date"],
		Phone:        a["phone"],
		Email:        a["email"],
		Organization: a["organization"],
		ConsentAt:    sess.ConsentAt,
	})
	repo
	switch {
	case err == nil:
		shell.sessions.Delete(from.ID)
		shell.deleteCallbackMessage(ctx, tgcl, update)
		shell.sendKey(ctx, tgcl, from.ID, messages.RegistrationDone, nil, nil)

	case errors.Is(err, user.ErrAlreadyExists):
		shell.sessions.Delete(from.ID)
		shell.deleteCallbackMessage(ctx, tgcl, update)
		shell.sendKey(ctx, tgcl, from.ID, messages.AlreadyRegistered, nil, nil)

	default:
		// даже если диск не доступен, тип может тыкнуть еще раз
		shell.logger.Error("Регистрация не сохранена", "telegram_id", from.ID, "error", err)
		shell.sendKey(ctx, tgcl, from.ID, messages.TryLater, nil, nil)
	}
}

// обработка измнения данных (Регистрация)
func (shell *Bot) editHandler(ctx context.Context, tgcl *tgbot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}
	from := update.CallbackQuery.From
	shell.answerCallback(ctx, tgcl, update, "")
	shell.deleteCallbackMessage(ctx, tgcl, update)

	shell.sessions.Restart(from.ID) // нужно добавить: Step = 0, Answers очищаются, согласие остаётся
	shell.sendKey(ctx, tgcl, from.ID, registration.Fields[0].Prompt, nil, nil)
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

	id := update.Message.From.ID

	switch res, f := shell.sessions.Submit(id, update.Message.Text); res {
	case registration.Invalid:
		shell.replyHandler(ctx, tgcl, update, f.Error, nil, nil)
	case registration.NextStep:
		shell.replyHandler(ctx, tgcl, update, f.Prompt, nil, nil)
	case registration.Finished:
		sess, ok := shell.sessions.Get(id)
		if !ok {
			return
		}
		shell.replyHandler(ctx, tgcl, update, messages.Summary, sess.Answers, confirmKeyboard())
	case registration.NoSession:
		shell.replyHandler(ctx, tgcl, update, messages.Unknown, nil, nil) // можете просто заменить на Help
	}
}
