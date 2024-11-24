package rabbit

import (
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
)

type Client struct {
	conf    *Config
	logg    *slog.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	Done    chan error
}

func New(
	logg *slog.Logger,
	conf Config,
	doneCh chan error,
) (*Client, error) {
	conn, err := amqp.Dial(conf.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conf:    &conf,
		logg:    logg,
		conn:    conn,
		channel: channel,
		Done:    doneCh,
	}, nil
}

func (c *Client) Close() error {
	err := c.channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	err = c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %w", err)
	}
	return nil
}

// RegisterOutboxExchange creates an exchange for outgoing messages.
func (c *Client) RegisterOutboxExchange() error {
	if err := c.channel.ExchangeDeclare(
		c.conf.ExchangeName,
		c.conf.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed exchange declare: %w", err)
	}

	return nil
}

// RegisterInboxQueue creates a queue for incoming messages.
func (c *Client) RegisterInboxQueue() error {
	queue, err := c.channel.QueueDeclare(
		c.conf.TopicName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %w", err)
	}

	if err = c.channel.QueueBind(
		queue.Name,
		c.conf.RoutingKey,
		c.conf.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("queue Bind: %w", err)
	}

	return nil
}
