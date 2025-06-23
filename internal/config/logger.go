package config

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func NewLogger(env string) *slog.Logger {
	var handler slog.Handler
	w := os.Stdout

	switch env {
	case devEnv:
		handler = tint.NewHandler(w, &tint.Options{Level: slog.LevelDebug})
	case prodEnv:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	case testEnv:
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	return slog.New(handler)
}
