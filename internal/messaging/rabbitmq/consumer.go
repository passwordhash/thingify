package rabbitmq

import (
	"context"
	"fmt"

	mq "github.com/rabbitmq/amqp091-go"
)

// сonsumer представляет собой потребителя RabbitMQ, который подключается к
// RabbitMQ-серверу, объявляет очередь и обрабатывает сообщения из нее.
type consumer struct {
	ch        *mq.Channel
	queueName string
}

// Consume начинает прослушивание очереди RabbitMQ и возвращает каналы
// для получения данных и ошибок. Данные будут отправляться в канал
// data, а ошибки в канал errCh.
func (c *consumer) Consume(ctx context.Context) (<-chan []byte, <-chan error) {
	const op = "rabbitmq.Consume"

	data := make(chan []byte)
	errCh := make(chan error, 100)
	go func() {
		defer close(data)
		defer close(errCh)

		msgs, err := c.ch.Consume(
			c.queueName, "", false, false, false, false, nil,
		)
		if err != nil {
			errCh <- fmt.Errorf("Consume: %w", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				if err := msg.Ack(false); err != nil {
					errCh <- fmt.Errorf("%s: ack failed: %w", op, err)
					return
				}
				data <- msg.Body
			}
		}
	}()

	return data, errCh
}

func (c *consumer) Close() error {
	return c.ch.Close()
}
