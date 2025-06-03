package server

import (
	"log/slog"
	"os"
	"thingify/internal/service/monitor"
	"time"

	"context"
)

type App struct {
	log *slog.Logger

	monitorService *monitor.Service
}

func New(
	log *slog.Logger,
	monitorService *monitor.Service,
) *App {
	return &App{
		log:            log,
		monitorService: monitorService,
	}
}

// MustRun запускает сервер и вызывает панику в случае ошибки.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic("failed to run server")
	}
}

// Run запускает сервер
func (a *App) Run() error {
	const op = "server.Run"

	log := a.log.With("op", op)

	log.Info("starting server")

	a.monitorService.ShortPollingNewIssues(
		context.TODO(), os.Getenv("GH_TOKEN"), 15*time.Second)

	return nil
}
