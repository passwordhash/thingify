package app

import (
	"context"
	"log/slog"
	httpapp "thingify/internal/app/http"
	"thingify/internal/config"
	issueService "thingify/internal/service/issue"
)

type App struct {
	Srv *httpapp.App
}

func New(
	_ context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	issueSvc := issueService.New(log)

	a := cfg.App
	h := cfg.HTTP
	srv := httpapp.New(
		log.WithGroup("http"),
		issueSvc,
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
