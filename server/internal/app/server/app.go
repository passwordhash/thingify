package server

import (
	"log/slog"
	"thingify/server/internal/github"
)

type App struct {
	log *slog.Logger

	ghClient *github.GHClient
}

func New(
    logger *slog.Logger,
    port int,
    ghBaseURL string,
    ghToken string,
) *App {

	ghClient := github.Register(logger, ghBaseURL, ghToken)

	return &App{
		log:      logger,
		ghClient: ghClient,
	}
}

func (a *App) MustRun() {

}
