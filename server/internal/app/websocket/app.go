package websocket

import (
	"log/slog"
	"net/http"
	"thingify/server/internal/github"
	"thingify/server/internal/websocket"
)

type App struct {
	log *slog.Logger

	hub      *websocket.Hub
	ghClient *github.GHClient
}

func New(
	logger *slog.Logger,
	port int,
	ghBaseURL string,
	ghToken string,
) *App {
	hub := websocket.NewHub()

	ghClient := github.Register(logger, ghBaseURL, ghToken)

	return &App{
		log:      logger,
		hub:      hub,
		ghClient: ghClient,
	}
}

func (a *App) MustRun() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", websocket.Handler(a.hub))

	addr := ":8080" // TODO: port from config
	a.log.Info("Starting server", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		a.log.Error("Server error", "err", err)
	}
}
