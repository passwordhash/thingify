package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"thingify/internal/app"
	"thingify/internal/config"

	"context"

	"github.com/gofiber/fiber/v2/log"
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

	logger := config.NewLogger(cfg.App.ENV)
	slog.SetDefault(logger) // TODO: нужно ли это?

	logger.Info("starting Thingify application http_server...")

	application := app.New(ctx, logger, cfg)

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
