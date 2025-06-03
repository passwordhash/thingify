package app

import (
	"log/slog"
	"thingify/internal/app/server"
	"thingify/internal/config"
	"thingify/internal/github"
	"thingify/internal/service/monitor"
)

type App struct {
	Server *server.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	ghClient := github.Register(log, cfg.GH.BaseURL, cfg.GHQueriesPath)

	monitorService := monitor.New(log, ghClient)

	srvApp := server.New(log, monitorService)

	return &App{
		Server: srvApp,
	}
}
