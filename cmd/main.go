package main

import (
	"log/slog"
	"os"
	"thingify/internal/app"
	"thingify/internal/config"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panic", "err", r)
		}
	}()

	w := os.Stdout
	log := slog.New(tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))

	cfg := config.MustLoad()

	application := app.New(log, cfg)

	application.Server.MustRun()
}
