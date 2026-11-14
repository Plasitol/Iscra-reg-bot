package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	// логирование
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// секреты + конфиг
	err := godotenv.Load()
	if err != nil {
		logger.Error(".env файл поврежден или не найден: ", err)
		os.Exit(1)
	}

	//
	father_key := os.Getenv("BOT_FATHER")
	//

	// оболочка бота
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, startHandler),
	}

	shellbot, err := bot.New(father_key, opts...)
	if err != nil {
		logger.Error("Бот не может запуститься...", err)
		client := &http.Client{Timeout: 3 * time.Second} // Проверка соединения
		resp, err := client.Get("https://telegram.org")
		if err != nil {
			logger.Error("Проверьте стабильность интернет соединения/ВПН для работы бота", err)
		}
		defer resp.Body.Close()
	}

	shellbot.Start(ctx)
}
