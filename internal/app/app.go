package app

import (
	"log/slog"
	"thingify/internal/app/server"
	"thingify/internal/config"
)

type App struct {
	Server *server.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	port := 8080 // TMP
	srvApp := server.New(log, port, cfg.GH.BaseURL, cfg.GHQueriesPath)

	return &App{
		Server: srvApp,
	}
}
