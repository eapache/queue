package queue

import "testing"

func TestQueueLength(t *testing.T) {
	q := New()

	if q.Length() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		q.Add(i)
		if q.Length() != i+1 {
			t.Error("adding: queue with", i, "elements has length", q.Length())
		}
	}
	for i := 0; i < 1000; i++ {
		q.Remove()
		if q.Length() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", q.Length())
		}
	}
}

func TestQueueGet(t *testing.T) {
	q := New()

	for i := 0; i < 1000; i++ {
		q.Add(i)
		for j := 0; j < q.Length(); j++ {
			if q.Get(j).(int) != j {
				t.Errorf("index %d doesn't contain %d", j, j)
			}
		}
	}
}

func TestQueueGetOutOfRangePanics(t *testing.T) {
	q := New()

	q.Add(1)
	q.Add(2)
	q.Add(3)

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("should panic when negative index")
			} else {
				t.Logf("got panic as expected: %v", r)
			}
		}()

		func() {
			q.Get(-1)
		}()
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("should panic when index greater than length")
			} else {
				t.Logf("got panic as expected: %v", r)
			}
		}()

		func() {
			q.Get(4)
		}()
	}()
}

// General warning: Go's benchmark utility (go test -bench .) increases the number of
// iterations until the benchmarks take a reasonable amount of time to run; memory usage
// is *NOT* considered. On my machine, these benchmarks hit around ~1GB before they've had
// enough, but if you have less than that available and start swapping, then all bets are off.

func BenchmarkQueueSerial(b *testing.B) {
	q := New()
	for i := 0; i < b.N; i++ {
		q.Add(nil)
	}
	for i := 0; i < b.N; i++ {
		q.Remove()
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := New()
	for i := 0; i < b.N; i++ {
		q.Add(nil)
		q.Remove()
	}
}
