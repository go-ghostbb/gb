package gbcron

import (
	"context"
	"ghostbb.io/gb/internal/intlog"
	"time"
)

// getAndUpdateLastTimestamp checks fixes and returns the last timestamp that have delay fix in some seconds.
func (s *cronSchedule) getAndUpdateLastTimestamp(ctx context.Context, t time.Time) int64 {
	var (
		currentTimestamp = t.Unix()
		lastTimestamp    = s.lastTimestamp.Val()
	)
	switch {
	case
		lastTimestamp == currentTimestamp:
		lastTimestamp += 1

	case
		lastTimestamp == currentTimestamp-1:
		lastTimestamp = currentTimestamp

	case
		lastTimestamp == currentTimestamp-2,
		lastTimestamp == currentTimestamp-3:
		lastTimestamp += 1

	default:
		// Too much delay, let's update the last timestamp to current one.
		intlog.Printf(
			ctx,
			`too much delay, last timestamp "%d", current "%d"`,
			lastTimestamp, currentTimestamp,
		)
		lastTimestamp = currentTimestamp
	}
	s.lastTimestamp.Set(lastTimestamp)
	return lastTimestamp
}
