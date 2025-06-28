package app

import (
	"context"
	"fmt"
	"log/slog"

	httpapp "thingify/internal/app/http"
	"thingify/internal/config"
	"thingify/internal/messaging/rabbitmq"
	issueService "thingify/internal/service/issue"
	"thingify/internal/storage/inmemory"
)

type App struct {
	Srv *httpapp.App
}

func New(
	_ context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	inmem := inmemory.New()

	rabbitmqClient, err := rabbitmq.NewClient(cfg.RabbitMQ.URL())
	if err != nil {
		panic(fmt.Errorf("failed to create RabbitMQ client: %w", err))
	}

	mqProducer, err := rabbitmqClient.NewProducer(cfg.RabbitMQ.IssueExchange)
	if err != nil {
		panic(fmt.Errorf("failed to create RabbitMQ producer: %w", err))
	}

	issueSvc := issueService.New(
		log.WithGroup("service"),
		mqProducer,
		inmem,
		inmem,
	)

	a := cfg.App
	h := cfg.HTTP
	srv := httpapp.New(
		log.WithGroup("http"),
		issueSvc,
		a.GithubWebhookSecret,
		httpapp.WithPort(h.Port),
		httpapp.WithReadTimeout(h.ReadTimeout),
		httpapp.WithWriteTimeout(h.ReadTimeout),
		httpapp.WithRequestTimeout(h.ReadTimeout),
	)

	return &App{
		Srv: srv,
	}
}
