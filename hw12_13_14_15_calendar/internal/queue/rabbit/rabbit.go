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
	Done    chan error
}

func New(
	logg *slog.Logger,
	conf Config,
	doneCh chan error,
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
		Done:    doneCh,
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
	queue, err := s.channel.QueueDeclare(
		s.conf.TopicName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	if err = s.channel.QueueBind(
		queue.Name,
		s.conf.RoutingKey,
		s.conf.ExchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

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
func (s *Service) Consume(messageCh chan string) error {
	deliveries, err := s.channel.Consume(
		s.conf.TopicName, // name
		"calendar-consumer",
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed queue Consume: %s", err)
	}

	go func() {
		for d := range deliveries {
			messageCh <- string(d.Body)
		}
		close(messageCh)
	}()

	return nil
}
