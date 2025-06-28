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
	ch       *mq.Channel
	exchange string
}

func (p *producer) Publish(ctx context.Context, routingKey string, payload any) error {
	const op = "rabbitmq.producer.Publish"

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal payload: %w", op, err)
	}

	// TODO: handle error
	p.ch.PublishWithContext(ctx,
		p.exchange,
		routingKey,
		false, false,
		mq.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	return nil
}

func (p *producer) Close() error {
	return p.ch.Close()
}
