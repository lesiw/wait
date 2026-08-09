package wait

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestDecayFirstIsFree(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var (
			w     = Decay(15 * time.Second)
			start = time.Now()
		)
		<-w.Wait()
		if got := time.Since(start); got != 0 {
			t.Errorf("Wait() first call blocked %v, want 0", got)
		}
	})
}

func TestDecayFreeUnderThreshold(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := Decay(15 * time.Second)
		for i := range 5 {
			start := time.Now()
			<-w.Wait()
			if got := time.Since(start); got != 0 {
				t.Errorf("Wait() call #%d blocked %v, want 0", i+1, got)
			}
		}
	})
}

func TestDecayCrossingThresholdBlocks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := Decay(15 * time.Second)
		for range 5 {
			<-w.Wait()
		}
		start := time.Now()
		<-w.Wait()
		if got, want := time.Since(start), 15*time.Second; got != want {
			t.Errorf("Wait() at threshold+1 blocked %v, want %v", got, want)
		}
	})
}

func TestDecayResetsAfterCrossing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := Decay(15 * time.Second)
		for range 6 {
			<-w.Wait()
		}
		start := time.Now()
		<-w.Wait()
		if got := time.Since(start); got != 0 {
			t.Errorf("Wait() after reset blocked %v, want 0", got)
		}
	})
}

func TestDecayHalfLifeReducesCount(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := Decay(15 * time.Second)
		for range 5 {
			<-w.Wait()
		}
		// Two half-lives (2 * 30s): 5 * 0.5^2 = 1.25.
		time.Sleep(60 * time.Second)
		for i := range 3 {
			start := time.Now()
			<-w.Wait()
			if got := time.Since(start); got != 0 {
				t.Errorf(
					"Wait() after decay, call #%d blocked %v, want 0",
					i+1, got,
				)
			}
		}
	})
}

func TestDecayReturnsClosedChannel(t *testing.T) {
	select {
	case <-Decay(15 * time.Second).Wait():
	default:
		t.Fatal("Wait() returned an unready channel; want pre-closed")
	}
}
