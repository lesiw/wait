package wait

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestJitterFirstWithinBase(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const base = 1 * time.Millisecond
		w := Jitter(1024 * time.Millisecond)
		start := time.Now()
		<-w.Wait()
		got := time.Since(start)
		if got < 0 || got >= base {
			t.Errorf(
				"Wait() first call blocked %v, want in [0, %v)", got, base,
			)
		}
	})
}

func TestJitterUpperBoundDoubles(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const base = 1 * time.Millisecond
		w := Jitter(1024 * time.Millisecond)
		bounds := []time.Duration{
			1 * base,
			2 * base,
			4 * base,
			8 * base,
			16 * base,
		}
		for i, upper := range bounds {
			start := time.Now()
			<-w.Wait()
			got := time.Since(start)
			if got < 0 || got >= upper {
				t.Errorf("Wait() call #%d blocked %v, want in [0, %v)",
					i+1, got, upper)
			}
		}
	})
}

func TestJitterSaturatesAtCap(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		limit := 1024 * time.Millisecond
		w := Jitter(limit)
		for range 20 {
			<-w.Wait()
		}
		for i := range 5 {
			start := time.Now()
			<-w.Wait()
			got := time.Since(start)
			if got < 0 || got >= limit {
				t.Errorf("saturated call #%d blocked %v, want in [0, %v)",
					i+1, got, limit)
			}
		}
	})
}

func TestJitterShiftOverflowClampsToCap(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		limit := 1024 * time.Second // internal base = 1s
		w := Jitter(limit)
		// Advance past shift overflow: base<<34 wraps to negative.
		for range 34 {
			<-w.Wait()
		}
		start := time.Now()
		<-w.Wait()
		if got := time.Since(start); got < 0 || got >= limit {
			t.Errorf(
				"Wait() past shift-overflow blocked %v, want in [0, %v)",
				got, limit,
			)
		}
	})
}

func TestJitterDistributionSpread(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const limit = 1024 * time.Millisecond
		const runs = 1000
		w := Jitter(limit)
		// Drain the pre-limit doublings so every draw is uniform on
		// [0, limit).
		for range 11 {
			<-w.Wait()
		}
		var (
			total time.Duration
			minS  = limit
			maxS  time.Duration
		)
		for range runs {
			start := time.Now()
			<-w.Wait()
			d := time.Since(start)
			total += d
			if d < minS {
				minS = d
			}
			if d > maxS {
				maxS = d
			}
		}
		mean := total / runs
		t.Logf("min=%v max=%v mean=%v (runs=%d)", minS, maxS, mean, runs)
		// Uniform on [0, limit): mean ~ limit/2.
		if mean < limit*40/100 || mean > limit*60/100 {
			t.Errorf("mean of %d draws = %v, want ~%v (uniform on [0, %v))",
				runs, mean, limit/2, limit)
		}
		if maxS-minS < limit/2 {
			t.Errorf("spread max-min = %v, want at least %v "+
				"(uniform on [0, %v))", maxS-minS, limit/2, limit)
		}
	})
}
