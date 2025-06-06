package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	mq "github.com/rabbitmq/amqp091-go"
)

const (
	directExchangeKind = "direct"
)

type producer struct {
	conn          *mq.Connection
	chann         *mq.Channel
	issueExchange string
}

func New(url, issueExchange string) (*producer, error) {
	const op = "rabbitmq.New"

	conn, err := mq.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to RabbitMQ: %w", op, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("%s: failed to open channel: %w", op, err)
	}

	err = ch.ExchangeDeclare(
		issueExchange,
		directExchangeKind,
		false,
		false, // autodelete
		false,
		false,
		nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("%s: failed to declare exchange %s: %w", op, issueExchange, err)
	}

	return &producer{
		conn:          conn,
		chann:         ch,
		issueExchange: issueExchange,
	}, nil
}

func (p *producer) Publish(ctx context.Context, routingKey string, payload any) error {
	const op = "rabbitmq.producer.Publish"

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal payload: %w", op, err)
	}

	const tmp_route_key = "users"

	p.chann.PublishWithContext(ctx,
		p.issueExchange,
		tmp_route_key,
		false, false,
		mq.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	return nil
}

func (p *producer) Close() error {
	_ = p.chann.Close()
	return p.conn.Close()
}
