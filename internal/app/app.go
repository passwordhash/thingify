package app

import (
	"context"
	"log/slog"
	"thingify/internal/app/server"
	"thingify/internal/config"
	"thingify/internal/github"
	"thingify/internal/service/monitor"
	"thingify/internal/storage/inmemory"
)

type App struct {
	Server *server.App
}

func New(_ context.Context, log *slog.Logger, cfg *config.Config) *App {
	storage := inmemory.New()

	ghClient := github.Register(log, cfg.GH.BaseURL, cfg.App.GHQueriesPath)

	monitorService := monitor.New(log, ghClient, storage, storage)

	srvApp := server.New(log, monitorService, cfg.App.PollingInterval)

	return &App{
		Server: srvApp,
	}
}
