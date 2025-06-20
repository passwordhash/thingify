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

	mqClient, err := rabbitmq.NewClient(cfg.RabbitMQ.URL())
	if err != nil {
		log.ErrorContext(ctx, "failed to create RabbitMQ client", slog.Any("error", err))
		os.Exit(1)
	}

	producer, err := mqClient.NewProducer(cfg.RabbitMQ.IssueExchange)
	if err != nil {
		log.ErrorContext(ctx, "failed to create RabbitMQ producer", slog.Any("error", err))
		os.Exit(1)
	}
	log.InfoContext(ctx, "RabbitMQ producer created successfully")

	consumer, err := mqClient.NewConsumer(cfg.RabbitMQ.CheckRequestsQueue)
	if err != nil {
		log.ErrorContext(ctx, "failed to create RabbitMQ consumer", slog.Any("error", err))
		os.Exit(1)
	}
	log.InfoContext(ctx, "RabbitMQ consumer created successfully")

	ghClient := github.Register(log, cfg.GH.BaseURL, cfg.App.GHQueriesPath)

	monitorService := monitor.New(log, ghClient, storage, storage, producer, consumer)

	srvApp := server.New(log, monitorService, cfg.App.PollingInterval)

	return &App{
		Server: srvApp,
	}
}
