package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	ghostshell "plasitol.dev/iscrav1/internal/bot"
)

func main() {
	// логирование
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("Пробуем запуститься...")

	// секреты + конфиг
	err := godotenv.Load()
	if err != nil {
		logger.Error(".env файл поврежден или не найден: ", "error", err)
		os.Exit(1)
	}

	father_key := os.Getenv("BOT_TOKEN")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM, // предохранитель
	)
	defer stop()

	// Инициализруем бота
	appBot, err := ghostshell.New(
		father_key,
		logger,
	)
	if err != nil {
		logger.Error("Бот не может запуститься... ", "error", err)
		client := &http.Client{Timeout: 3 * time.Second} // Проверка соединения
		resp, err := client.Get("https://telegram.org")
		if err != nil {
			logger.Error("Проверьте стабильность интернет соединения/ВПН для работы бота: ", "error", err)
		}
		defer resp.Body.Close()
		os.Exit(1)
	}

	logger.Info("Бот запущен!")
	appBot.Start(ctx)
}
