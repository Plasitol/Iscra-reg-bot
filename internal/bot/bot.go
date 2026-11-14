package bot

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram/bot"
	"plasitol.dev/iscrav1/internal/messages"
	"plasitol.dev/iscrav1/internal/user"
)

type Bot struct {
	client *tgbot.Bot
	logger *slog.Logger
	msgs   *messages.Store
	users  *user.Service
}

func New(token string, logger *slog.Logger, msgs *messages.Store, users *user.Service) (*Bot, error) {
	shell := &Bot{
		logger: logger,
		msgs:   msgs,
		users:  users,
	}

	client, err := tgbot.New(token, shell.routes()...)
	if err != nil {
		return nil, err
	}

	shell.client = client

	return shell, nil
}

func (shell *Bot) Start(ctx context.Context) {
	shell.client.Start(ctx)
}
