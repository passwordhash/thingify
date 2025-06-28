package rabbitmq

import (
	"fmt"

	mq "github.com/rabbitmq/amqp091-go"
)

const (
	IssueExchangeName = "github_issue" // TODO: подумать о целесообразности
)

type Client struct {
	conn *mq.Connection
}

// NewClient создает нового клиента RabbitMQ, устанавливает соединение.
func NewClient(url string) (*Client, error) {
	const op = "messaging.rabbitmq.NewClient"

	conn, err := mq.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to RabbitMQ: %w", op, err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) NewProducer(exchange string) (*producer, error) {
	const op = "messaging.rabbitmq.NewProducer"

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open channel: %w", op, err)
	}

	err = ch.ExchangeDeclare(
		exchange,
		directExchangeKind,
		false,
		false, // autodelete
		false,
		false,
		nil)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("%s: failed to declare exchange %s: %w", op, exchange, err)
	}

	return &producer{
		ch:       ch,
		exchange: exchange,
	}, nil
}

func (c *Client) NewConsumer(queueName, routingKey, issueExchange string) (*consumer, error) {
	const op = "messaging.rabbitmq.NewConsumer"

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open new channel: %w", op, err)
	}

	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to declare queue %s: %w", op, queueName, err)
	}

	err = ch.QueueBind(queueName, routingKey, issueExchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to bind queue %s with routing key %s to exchange %s: %w",
			op, queueName, routingKey, issueExchange, err,
		)
	}

	return &consumer{
		ch:        ch,
		queueName: queue.Name,
	}, nil

}

func (c *Client) Close() error {
	return c.conn.Close()
}
