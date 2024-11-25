package scheduler

import (
	"context"
	"time"
)

func (a *AppScheduler) cleanupEvents(ctx context.Context) func(d time.Duration) {
	return func(t time.Duration) {
		a.logg.Debug("Cleanup task started")

		pastTime := time.Now().Add(-t)
		err := a.repo.DeleteOlderThan(ctx, pastTime)
		if err != nil {
			a.logg.Error("failed to clean events: " + err.Error())
			return
		}
	}
}
