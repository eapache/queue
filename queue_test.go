package queue

import "testing"

func TestQueueSimple(t *testing.T) {
	q := New()

	for i := 0; i < minQueueLen; i++ {
		q.Add(i)
	}
	for i := 0; i < minQueueLen; i++ {
		if e, _ := q.Peek(); e.(int) != i {
			t.Error("peek", i, "had value", e)
		}
		q.Remove()
	}
}

func TestQueueWrapping(t *testing.T) {
	q := New()

	for i := 0; i < minQueueLen; i++ {
		q.Add(i)
	}
	for i := 0; i < 3; i++ {
		q.Remove()
		q.Add(minQueueLen + i)
	}

	for i := 0; i < minQueueLen; i++ {
		if e, _ := q.Peek(); e.(int) != i+3 {
			t.Error("peek", i, "had value", e)
		}
		q.Remove()
	}
}

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
			if e, _ := q.Get(j); e.(int) != j {
				t.Errorf("index %d doesn't contain %d", j, j)
			}
		}
	}
}

func TestQueuePops(t *testing.T) {
	q := New()

	for i := 0; i < 1000; i++ {
		q.Add(i)
	}

	for i := 0; i < 1000; i++ {
		if e, _ := q.Pop(); e != i {
			t.Errorf("index %d doesn't contain %d", i, i)
		}
	}
}

func TestQueueGetOutOfRangeErrors(t *testing.T) {
	q := New()

	q.Add(1)
	q.Add(2)
	q.Add(3)

	_, err := q.Get(-1)
	if err == nil {
		t.Error("should haved errored when negative index")
	}

	_, err = q.Get(4)
	if err == nil {
		t.Error("should haved errored when negative index")
	}
}

func TestQueuePeekOutOfRangeErrors(t *testing.T) {
	q := New()

	if _, err := q.Peek(); err == nil {
		t.Error("should error when peeking empty queue")
	}

	q.Add(1)
	q.Remove()

	if _, err := q.Peek(); err == nil {
		t.Error("should error when peeking emptied queue")
	}
}

func TestQueueRemoveOutOfRangeErrors(t *testing.T) {
	q := New()

	if q.Remove() == nil {
		t.Error("should error when removing empty queue")
	}

	q.Add(1)
	q.Remove()

	if q.Remove() == nil {
		t.Error("should error when removing emptied queue")
	}
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
		q.Peek()
		q.Remove()
	}
}

func BenchmarkQueueGet(b *testing.B) {
	q := New()
	for i := 0; i < b.N; i++ {
		q.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Get(i)
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := New()
	for i := 0; i < b.N; i++ {
		q.Add(nil)
		q.Peek()
		q.Remove()
	}
}
