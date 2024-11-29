package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Publish sends a message to the queue.
func (c *Client) Publish(routingKey string, body []byte) error {
	if err := c.channel.Publish(
		c.conf.ExchangeName, // Вынести в параметры функции
		routingKey,          // Вынести в параметры функции
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,
		},
	); err != nil {
		return fmt.Errorf("failed to publish into exchange: %w", err)
	}
	return nil
}
