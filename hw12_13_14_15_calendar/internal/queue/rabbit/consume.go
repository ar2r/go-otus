package rabbit

import (
	"fmt"
)

const ConsumerName = "calendar-consumer"

// Consume returns a channel with messages from the queue
func (c *Client) Consume(messageCh chan<- string) error {
	deliveries, err := c.channel.Consume(
		c.conf.TopicName, // name
		ConsumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed queue Consume: %c", err)
	}

	go func() {
		for d := range deliveries {
			messageCh <- string(d.Body)
			d.Ack(true)
		}
		close(messageCh)
	}()

	return nil
}
