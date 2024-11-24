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
	// Publish sends a message to the queue
	Publish(routingKey string, body []byte) error
}

type IConsumer interface {
	// Consume returns a channel with messages from the queue
	Consume() (<-chan []byte, error)
}

func NewProducer(
	logg *slog.Logger,
	conf *config.Config,
) (IProducer, error) {
	if conf.App.Queue == MessageBrokerRabbitMQ {
		rabbitConn, err := rabbit.New(logg, conf.RabbitMQ)
		if err != nil {
			return nil, fmt.Errorf("failed to create rabbitmq producer: %w", err)
		}

		// Register exchange
		if err := rabbitConn.RegisterOutboxExchange(); err != nil {
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

//func NewConsumer(ctx context.Context, logg *slog.Logger, conf *Config) (IConsumer, error) {
//	return nil, nil
//}
