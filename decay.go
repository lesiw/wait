package wait

import (
	"math"
	"time"
)

// Decay returns a waiter whose first five calls are free. The sixth call
// blocks for limit and the count then resets. The count halves every 2*limit
// of idle time. Good for supervisor restart loops. Inspired by Suture's
// failure-decay counter (github.com/thejerf/suture/v4).
func Decay(limit time.Duration) Waiter {
	return &decayWaiter{limit: limit}
}

type decayWaiter struct {
	limit    time.Duration
	count    float64
	lastCall time.Time
}

func (w *decayWaiter) Wait() <-chan time.Time {
	now := time.Now()
	if !w.lastCall.IsZero() {
		intervals := float64(now.Sub(w.lastCall)) / float64(2*w.limit)
		w.count *= math.Exp2(-intervals)
	}
	w.count++
	w.lastCall = now
	if w.count > 5 {
		w.count = 0
		return time.After(w.limit)
	}
	return closed
}
