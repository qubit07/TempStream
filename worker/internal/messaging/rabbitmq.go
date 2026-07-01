package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	deliveries <-chan amqp.Delivery
}

type Config struct {
	URL      string
	Queue    string
	Prefetch int
}

func NewConsumer(cfg Config) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	queue, err := channel.QueueDeclare(
		cfg.Queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("queue declare: %w", err)
	}

	if cfg.Prefetch > 0 {
		if err := channel.Qos(cfg.Prefetch, 0, false); err != nil {
			_ = channel.Close()
			_ = conn.Close()
			return nil, fmt.Errorf("qos: %w", err)
		}
	}

	deliveries, err := channel.Consume(
		queue.Name,
		"",
		false,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &Consumer{
		conn:       conn,
		channel:    channel,
		queue:      queue,
		deliveries: deliveries,
	}, nil
}

func (c *Consumer) Run(
	ctx context.Context,
	handler func(context.Context, amqp.Delivery) error,
) error {

	log.Printf("consuming queue: %s", c.queue.Name)

	for {
		select {

		case <-ctx.Done():
			log.Println("consumer shutdown requested")
			return ctx.Err()

		case msg, ok := <-c.deliveries:
			if !ok {
				return fmt.Errorf("deliveries channel closed")
			}

			err := handler(ctx, msg)
				if err != nil {
					log.Printf("handler error: %v", err)
				}
		}
	}
}

func (c *Consumer) Close() error {
	log.Println("closing rabbitmq consumer")

	if c.channel != nil {
		_ = c.channel.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}