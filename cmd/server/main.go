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

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panic", "err", r)
		}
	}()

	ctx := context.Background()

	cfg := config.MustLoad()

	w := os.Stdout
	log := slog.New(tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))

	log.Info("starting Thingify application server...")
	log.Debug("with config", "config", cfg)

	application := app.New(ctx, log, cfg)

	go application.Server.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("received signal", "signal", sign)

	application.Server.Stop()

	log.Info("stopped Thingify application server")
}
