package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"voting-bot/internal/app"
	"voting-bot/internal/config"
)

func main() {

	cfg := config.LoadConfig()

	log := setupLogger()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application := app.NewApp(log, cfg)
	go func() {
		if err := application.Run(); err != nil {
			log.Error(fmt.Errorf("error running server %w", err).Error())
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down gracefully...")

	application.Stop()

	log.Info("Bot stopped")
}

func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}
