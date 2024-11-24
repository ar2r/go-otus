package rabbit

import (
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
)

type Service struct {
	conf    *Config
	logg    *slog.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(
	logg *slog.Logger,
	conf Config,
) (*Service, error) {
	conn, err := amqp.Dial(conf.Uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %s", err)
	}

	return &Service{
		conf:    &conf,
		logg:    logg,
		conn:    conn,
		channel: channel,
	}, nil
}

func (s *Service) Close() error {
	err := s.channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close channel: %s", err)
	}

	err = s.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %s", err)
	}
	return nil
}

// RegisterOutboxExchange creates an exchange for outgoing messages
func (s *Service) RegisterOutboxExchange() error {
	if err := s.channel.ExchangeDeclare(
		s.conf.ExchangeName,
		s.conf.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed exchange declare: %s", err)
	}

	return nil
}

// RegisterInboxQueue creates a queue for incoming messages
func (s *Service) RegisterInboxQueue() error {
	return nil
}

// Publish sends a message to the queue
// "event.notification.upcoming"
func (s *Service) Publish(routingKey string, body []byte) error {
	if err := s.channel.Publish(
		s.conf.ExchangeName, // Вынести в параметры функции
		routingKey,          // Вынести в параметры функции
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("failed to publish into exchange: %s", err)
	}
	return nil
}

// Consume returns a channel with messages from the queue
func (s *Service) Consume() (<-chan []byte, error) {
	return nil, nil
}
