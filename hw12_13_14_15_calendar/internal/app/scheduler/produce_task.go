package scheduler

import (
	"context"
	"encoding/json"
)

func (a *AppScheduler) produceNotification(ctx context.Context) func(queueConn MessageProducer, routingKey string) {
	return func(queueConn MessageProducer, routingKey string) {
		a.logg.Debug("Produce notification task started")

		events, err := a.repo.ListNotNotified(ctx)
		if err != nil {
			a.logg.Error("failed to get events: " + err.Error())
			a.errorCh <- err
			return
		}
		for _, event := range events {
			eventJSON, err := json.Marshal(event)
			if err != nil {
				a.logg.Warn("failed to marshal event: " + err.Error())
				continue
			}
			if err = queueConn.Publish(routingKey, eventJSON); err != nil {
				a.logg.Error("failed to publish message: " + err.Error())
				a.errorCh <- err
				return
			}
			event.NotificationSent = true

			_, err = a.repo.Update(ctx, event)
			if err != nil {
				a.logg.Error("failed to update event: " + err.Error())
				a.errorCh <- err
			}
			a.logg.Debug("Event has been sent to MQ: " + event.ID.String())
		}
	}
}
