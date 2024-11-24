package kafka

// Пример конфига для Kafka. Просто для демонстрации того, что конфиги для разных брокеров очередей могут отличаться.

type Config struct {
	Hosts string `toml:"hosts"`
}
