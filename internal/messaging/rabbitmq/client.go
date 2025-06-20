package rabbitmq

import (
	"fmt"

	mq "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *mq.Connection
}

// NewClient создает нового клиента RabbitMQ, устанавливает соединение.
func NewClient(url string) (*Client, error) {
	const op = "rabbitmq.NewClient"

	conn, err := mq.Dial(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) NewProducer(issueExchange string) (*producer, error) {
	const op = "rabbitmq.NewProducer"

	ch, err := c.conn.Channel()
	if err != nil {
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
		return nil, fmt.Errorf("%s: failed to declare exchange %s: %w", op, issueExchange, err)
	}

	return &producer{
		ch:            ch,
		issueExchange: issueExchange,
	}, nil
}

func (c *Client) NewConsumer(queueName string) (*consumer, error) {
	const op = "rabbitmq.NewConsumer"

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open new channel: %w", op, err)
	}

	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to declare queue %s: %w", op, queueName, err)
	}

	return &consumer{
		ch:        ch,
		queueName: queue.Name,
	}, nil

}

func (c *Client) Close() error {
	return c.conn.Close()
}
