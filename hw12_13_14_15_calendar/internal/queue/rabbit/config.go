package rabbit

type Config struct {
	URI string `toml:"uri"`

	ExchangeName string `toml:"exchange_name"`
	ExchangeType string `toml:"exchange_type"`

	TopicName  string `toml:"topic_name"`
	RoutingKey string `toml:"routing_key"`
}
