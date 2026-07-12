// Package wait provides waiting strategies. A [Waiter] returns a channel that
// fires when the caller may proceed.
package wait

import "time"

// Waiter returns a channel that fires after a strategy-determined interval. A
// Waiter is not safe for concurrent use.
type Waiter interface {
	// Wait returns a channel that fires when the caller may proceed.
	Wait() <-chan time.Time
}

var closed = make(chan time.Time)

func init() {
	close(closed)
}
