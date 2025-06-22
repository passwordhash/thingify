package app

import (
	"context"
	"log/slog"
	httpapp "thingify/internal/app/http"
	"thingify/internal/config"
)

type App struct {
	Srv *httpapp.App
}

func New(_ context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	a := cfg.App
	h := cfg.HTTP
	srv := httpapp.New(
		log,
		a.GithubWebhookSecret,
		httpapp.WithPort(h.Port),
		httpapp.WithReadTimeout(h.ReadTimeout),
		httpapp.WithWriteTimeout(h.ReadTimeout),
		httpapp.WithRequestTimeout(h.ReadTimeout),
	)

	return &App{
		Srv: srv,
	}
}
