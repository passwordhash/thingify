package main

import (
	"log/slog"
	"os"
	"thingify/server/internal/app"
	"thingify/server/internal/config"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	w := os.Stdout
	log := slog.New(tint.NewHandler(w, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
	}))

	cfg := config.MustLoad()

	application := app.New(log, cfg)

	application.Server.MustRun()
}
