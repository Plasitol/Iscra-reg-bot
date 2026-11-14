package bot

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram/bot"
	"plasitol.dev/iscrav1/internal/messages"
	registration "plasitol.dev/iscrav1/internal/registation"
	"plasitol.dev/iscrav1/internal/user"
)

type Bot struct {
	client   *tgbot.Bot
	logger   *slog.Logger
	msgs     *messages.Store
	users    *user.Service
	sessions *registration.Sessions
}

func New(token string, logger *slog.Logger, msgs *messages.Store, users *user.Service, sessions *registration.Sessions) (*Bot, error) {
	shell := &Bot{
		logger:   logger,
		msgs:     msgs,
		users:    users,
		sessions: sessions,
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
