package bot

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram/bot"
	"plasitol.dev/iscrav1/internal/messages"
)

type Bot struct {
	client *tgbot.Bot
	logger *slog.Logger
	msgs   *messages.Store
}

func New(token string, logger *slog.Logger, msgs *messages.Store) (*Bot, error) {
	shell := &Bot{
		logger: logger,
		msgs:   msgs,
	}
	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(shell.defaultHandler),
		tgbot.WithMessageTextHandler(
			"/start",
			tgbot.MatchTypeExact,
			shell.startHandler,
		),
		tgbot.WithMessageTextHandler(
			"/help",
			tgbot.MatchTypeExact,
			shell.helpHandler,
		),
	}

	client, err := tgbot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	shell.client = client

	return shell, nil
}

func (shell *Bot) Start(ctx context.Context) {
	shell.client.Start(ctx)
}
