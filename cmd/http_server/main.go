package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"thingify/internal/app"
	"thingify/internal/config"
	"time"

	"context"

	"github.com/lmittmann/tint"
)

const (
	shutdownTimeout = 10 * time.Second
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panic", "err", r)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	// TODO: вынести в отдельный пакет
	w := os.Stdout
	log := slog.New(tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))

	log.Info("starting Thingify application http_server...")
	log.Debug("with config", "config", cfg)

	application := app.New(ctx, log, cfg)

	go application.Srv.MustRun(ctx)

	<-ctx.Done()

	log.Info("received stop signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := application.Srv.Stop(shutdownCtx); err != nil {
		log.Error("failed to stop http_server gracefully", "err", err)
	} else {
		log.Info("Thingify application http_server stopped gracefully")
	}
}
