package server

import (
	"log/slog"
	"os"
	"thingify/internal/github"
)

type App struct {
	log *slog.Logger

	ghClient *github.GHClient
}

func New(
	logger *slog.Logger,
	port int,
	ghBaseURL string,
	ghQueriesPath string,
) *App {

	ghClient := github.Register(logger, ghBaseURL, ghQueriesPath)

	return &App{
		log:      logger,
		ghClient: ghClient,
	}
}

func (a *App) MustRun() {
	a.ghClient.UserIssues(os.Getenv("GH_TOKEN"))
}
