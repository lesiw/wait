package wait

import (
	"math/rand/v2"
	"time"
)

// Jitter returns a full-jitter exponential waiter. Successive Wait calls
// sleep for a random duration in [0, upper), where upper doubles each call
// starting at limit/1024 and saturating at limit. Default choice for retry
// pacing across concurrent peers.
func Jitter(limit time.Duration) Waiter {
	return &jitterWaiter{limit: limit}
}

type jitterWaiter struct {
	limit   time.Duration
	attempt int
}

func (w *jitterWaiter) Wait() <-chan time.Time {
	d := (w.limit / 1024) << w.attempt
	if d <= 0 {
		d = w.limit
	}
	d = min(d, w.limit)
	if d < w.limit {
		w.attempt++
	}
	return time.After(rand.N(d))
}
