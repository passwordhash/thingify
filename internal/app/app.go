package app

import (
	"context"
	"log/slog"
	"os"
	"thingify/internal/app/server"
	"thingify/internal/config"
	"thingify/internal/github"
	"thingify/internal/messaging/rabbitmq"
	"thingify/internal/service/monitor"
	"thingify/internal/storage/inmemory"
)

type App struct {
	Server *server.App
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) *App {
	storage := inmemory.New()

	producer, err := rabbitmq.New(cfg.RabbitMQ.URL(), cfg.RabbitMQ.IssueExchange)
	if err != nil {
		log.ErrorContext(ctx, "failed to create RabbitMQ producer", slog.Any("error", err))
		os.Exit(1)
	}

	ghClient := github.Register(log, cfg.GH.BaseURL, cfg.App.GHQueriesPath)

	monitorService := monitor.New(log, ghClient, storage, storage, producer)

	srvApp := server.New(log, monitorService, cfg.App.PollingInterval)

	return &App{
		Server: srvApp,
	}
}
