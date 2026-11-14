package bot

import (
	tgbot "github.com/go-telegram/bot"
)

// Тексты кнопок
const (
	btnRegistration = "Регистрация"
	btnHelp         = "Помощь"
)

// Callback данные инлайн кнопок
const (
	// обработка пд
	cbConsentAccept  = "consent:accept"
	cbConsentDecline = "consent:decline"

	// регистрация
	reRegistrationConfirm = "registration:confirm"
	reRegistraionEdit     = "registration:edit"
)

func (shell *Bot) routes() []tgbot.Option {
	return []tgbot.Option{
		tgbot.WithDefaultHandler(shell.defaultHandler),

		// команды
		tgbot.WithMessageTextHandler("/start", tgbot.MatchTypeExact, shell.startHandler),
		// tgbot.WithMessageTextHandler("/help", tgbot.MatchTypeExact, shell.helpHandler),

		// кнопки главного меню
		tgbot.WithMessageTextHandler(btnHelp, tgbot.MatchTypeExact, shell.helpHandler),
		tgbot.WithMessageTextHandler(btnRegistration, tgbot.MatchTypeExact, shell.registrationHandler),

		// инлайн кнопки
		// обработка пд
		tgbot.WithCallbackQueryDataHandler(cbConsentAccept, tgbot.MatchTypeExact, shell.consentAcceptHandler),
		tgbot.WithCallbackQueryDataHandler(cbConsentDecline, tgbot.MatchTypeExact, shell.consentDeclineHandler),
		// обработка регистрации
		tgbot.WithCallbackQueryDataHandler("registration:confirm", tgbot.MatchTypeExact, shell.confirmHandler),
		tgbot.WithCallbackQueryDataHandler("registration:edit", tgbot.MatchTypeExact, shell.editHandler),
	}
}
