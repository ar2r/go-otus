package queue

import (
	"fmt"
	"log/slog"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue/rabbit"
)

const (
	MessageBrokerRabbitMQ = "rabbitmq"
	MessageBrokerKafka    = "kafka"
)

var (
	ErrInvalidQueueType        = fmt.Errorf("invalid queue type")
	ErrQueueTypeNotImplemented = fmt.Errorf("queue type not implemented")
)

type IProducer interface {
	Publish(routingKey string, body []byte) error
	Close() error
}

type IConsumer interface {
	Consume(chan string) error
	Close() error
}

func NewProducer(
	logg *slog.Logger,
	conf *config.Config,
) (IProducer, error) {
	if conf.App.Queue == MessageBrokerRabbitMQ {
		rabbitConn, err := rabbit.New(logg, conf.RabbitMQ, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create rabbitmq producer: %w", err)
		}

		if err = rabbitConn.RegisterOutboxExchange(); err != nil {
			return nil, fmt.Errorf("failed to register outbox exchange: %w", err)
		}
		return rabbitConn, nil
	}

	if conf.App.Queue == MessageBrokerKafka {
		logg.Warn("Kafka producer is not implemented")
		// kafka.Now(logg, conf.Kafka
		return nil, ErrQueueTypeNotImplemented
	}

	return nil, ErrInvalidQueueType
}

func NewConsumer(
	logg *slog.Logger,
	conf *config.Config,
	doneCh chan error,
) (IConsumer, error) {
	if conf.App.Queue == MessageBrokerRabbitMQ {
		rabbitConn, err := rabbit.New(logg, conf.RabbitMQ, doneCh)
		if err != nil {
			return nil, fmt.Errorf("failed to create rabbitmq consumer: %w", err)
		}

		err = rabbitConn.RegisterInboxQueue()
		if err != nil {
			return nil, fmt.Errorf("failed to register inbox queue: %w", err)
		}

		return rabbitConn, nil
	}
	return nil, nil
}
