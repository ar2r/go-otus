package scheduler

import (
	"context"
	"encoding/json"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/queue"
)

func (a *AppScheduler) produceNotification(ctx context.Context) func(queueConn queue.IProducer, routingKey string) {
	return func(queueConn queue.IProducer, routingKey string) {
		a.logg.Debug("Produce notification task started")

		events, err := a.repo.ListNotNotified(ctx)
		if err != nil {
			a.logg.Error("failed to get events: " + err.Error())
			return
		}
		for _, event := range events {
			eventJSON, err := json.Marshal(event)
			if err != nil {
				a.logg.Error("failed to marshal event: " + err.Error())
				continue
			}
			if err = queueConn.Publish(routingKey, eventJSON); err != nil {
				a.logg.Error("failed to publish message: " + err.Error())
			}
			event.NotificationSent = true

			_, err = a.repo.Update(ctx, event)
			if err != nil {
				a.logg.Error("failed to update event: " + err.Error())
			}
			a.logg.Debug("Event has been sent to MQ: " + event.ID.String())
		}
	}
}
