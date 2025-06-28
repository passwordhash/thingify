package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"thingify/internal/config"
	"thingify/internal/messaging/rabbitmq"

	"github.com/lmittmann/tint"
)

const (
	issuesQueue = "issues.opened"
)

var rcfg = config.RabbitMQConfig{
	Host: "localhost",
	Port: 5672,
	User: "guest",
	Pass: "guest",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	installID := os.Getenv("INSTALL_ID")
	if _, err := strconv.Atoi(installID); err != nil {
		panic("INSTALL_ID must be a valid integer")
	}

	rabbitClient, err := rabbitmq.NewClient(rcfg.URL())
	if err != nil {
		panic(fmt.Errorf("failed to connect to RabbitMQ: %w", err))
	}
	defer func() { _ = rabbitClient.Close() }()

	slog.Info("Connected to RabbitMQ", "url", rcfg.URL())

	exchange := rabbitmq.IssueExchangeName // TODO: подумать о целесообразности
	consumer, err := rabbitClient.NewConsumer(issuesQueue, installID, exchange)
	if err != nil {
		panic(fmt.Errorf("failed to create consumer: %w", err))
	}
	defer func() { _ = consumer.Close() }()

	dataCh, errCh := consumer.Consume(ctx)

	slog.Info("Consumer started", "queue", issuesQueue)

	out, _ := os.Create("out.json")
	go func() {
		for data := range dataCh {
			_, err = out.Write(data)
			if err != nil {
				fmt.Printf("Error marshaling data: %v\n", err)
				continue
			}
		}
	}()

	go func() {
		for err := range errCh {
			slog.Error("Error received from RabbitMQ", "error", err)
		}
	}()

	<-ctx.Done()
}
