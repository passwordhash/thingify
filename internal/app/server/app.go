package server

import (
	"log/slog"
	"thingify/internal/service/monitor"
	"time"

	"context"
)

type App struct {
	log *slog.Logger

	monitorService *monitor.Service

	pollingInterval time.Duration
}

func New(
	log *slog.Logger,
	monitorService *monitor.Service,
	pollingInterval time.Duration,
) *App {
	return &App{
		log:             log,
		monitorService:  monitorService,
		pollingInterval: pollingInterval,
	}
}

// MustRun запускает сервер и вызывает панику в случае ошибки.
func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic("failed to run server")
	}
}

// Run запускает сервер
func (a *App) Run(ctx context.Context) error {
	const op = "server.Run"

	log := a.log.With("op", op)

	log.Info("starting server")

	a.monitorService.StartRequestListener(ctx)

	// a.monitorService.ShortPollingNewIssues(
	// 	ctx, os.Getenv("GH_TOKEN"), a.pollingInterval)

	return nil
}

// Stop останавливает сервер
// TODO: stop rabbitmq client, producer, consumer, etc.
func (a *App) Stop() {
	const op = "server.Stop"

	log := a.log.With("op", op)

	log.Info("stopping server")
}
