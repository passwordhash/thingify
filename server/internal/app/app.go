package app

import (
	"log/slog"
	"thingify/server/internal/app/server"
	"thingify/server/internal/config"
)

type App struct {
	Server *server.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	port := 8080 // TMP
	srvApp := server.New(log, port, cfg.GH.BaseURL, cfg.GH.Token)

	return &App{
		Server: srvApp,
	}
}
