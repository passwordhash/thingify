package app

import (
	"log/slog"
	"thingify/server/internal/app/websocket"
	"thingify/server/internal/config"
)

type App struct {
	Server *websocket.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	port := 8080 // TMP
	srvApp := websocket.New(log, port, cfg.GH.BaseURL, cfg.GH.Token)

	return &App{
		Server: srvApp,
	}
}
