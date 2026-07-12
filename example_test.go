package wait_test

import (
	"context"
	"fmt"
	"time"

	"lesiw.io/wait"
)

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var (
		start = time.Now()
		since = start
		tries = 0
	)
	fn := func() error {
		tries++
		now := time.Now()
		fmt.Printf("attempt=%d elapsed=%v delay=%v\n",
			tries,
			now.Sub(start).Round(time.Millisecond),
			now.Sub(since).Round(time.Millisecond),
		)
		since = now
		if tries < 5 {
			return fmt.Errorf("bang!")
		}
		return nil
	}
	wtr := wait.Jitter(2 * time.Second)
	for fn() != nil {
		select {
		case <-wtr.Wait():
		case <-ctx.Done():
			return
		}
	}
	fmt.Println("succeeded after", tries, "tries")
}
