package queue

import "testing"

func TestQueueSimple(t *testing.T) {
	q := New[int]()

	for i := 0; i < minQueueLen; i++ {
		q.Add(i)
	}
	for i := 0; i < minQueueLen; i++ {
		if r, _ := q.Peek(); r != i {
			t.Error("peek", i, "had value", r)
		}
		x, _ := q.Remove()
		if x != i {
			t.Error("remove", i, "had value", x)
		}
	}
}

func TestQueueWrapping(t *testing.T) {
	q := New[int]()

	for i := 0; i < minQueueLen; i++ {
		q.Add(i)
	}
	for i := 0; i < 3; i++ {
		_, _ = q.Remove()
		q.Add(minQueueLen + i)
	}

	for i := 0; i < minQueueLen; i++ {
		if r, _ := q.Peek(); r != i+3 {
			t.Error("peek", i, "had value", r)
		}
		_, _ = q.Remove()
	}
}

func TestQueueLength(t *testing.T) {
	q := New[int]()

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
		_, _ = q.Remove()
		if q.Length() != 1000-i-1 {
			t.Error("removing: queue with", 1000-i-i, "elements has length", q.Length())
		}
	}
}

func TestQueueGet(t *testing.T) {
	q := New[int]()

	for i := 0; i < 1000; i++ {
		q.Add(i)
		for j := 0; j < q.Length(); j++ {
			if r, _ := q.Get(j); r != j {
				t.Errorf("index %d doesn't contain %d", j, j)
			}
		}
	}
}

func TestQueueGetNegative(t *testing.T) {
	q := New[int]()

	for i := 0; i < 1000; i++ {
		q.Add(i)
		for j := 1; j <= q.Length(); j++ {
			if r, _ := q.Get(-j); r != q.Length()-j {
				t.Errorf("index %d doesn't contain %d", -j, q.Length()-j)
			}
		}
	}
}

func TestQueueGetOutOfRangePanics(t *testing.T) {
	q := New[int]()

	q.Add(1)
	q.Add(2)
	q.Add(3)

	_, err := q.Get(-4)
	assertError(t, "should get index out of range when negative index", err, ErrIndexOutOfRange)

	_, err = q.Get(4)
	assertError(t, "should panic when index greater than length", err, ErrIndexOutOfRange)
}

func TestQueuePeekOutOfRangePanics(t *testing.T) {
	q := New[any]()

	_, err := q.Peek()
	assertError(t, "should return empty queue error when peeking empty queue", err, ErrQueueEmpty)

	q.Add(1)
	_, _ = q.Remove()

	_, err = q.Peek()
	assertError(t, "should return empty queue error when peeking emptied queue", err, ErrQueueEmpty)
}

func TestQueueRemoveOutOfRangePanics(t *testing.T) {
	q := New[int]()

	_, err := q.Remove()
	assertError(t, "should return empty queue error when removing empty queue", err, ErrQueueEmpty)

	q.Add(1)
	_, _ = q.Remove()

	_, err = q.Remove()
	assertError(t, "should return empty queue error when removing emptied queue", err, ErrQueueEmpty)
}

func assertError(t *testing.T, name string, actualErr error, expectedErr error) {
	if actualErr != expectedErr {
		t.Errorf("%s: didn't get error as expected", name)
	}
}

// WARNING: Go's benchmark utility (go test -bench .) increases the number of
// iterations until the benchmarks take a reasonable amount of time to run; memory usage
// is *NOT* considered. On a fast CPU, these benchmarks can fill hundreds of GB of memory
// (and then hang when they start to swap). You can manually control the number of iterations
// with the `-benchtime` argument. Passing `-benchtime 1000000x` seems to be about right.

func BenchmarkQueueSerial(b *testing.B) {
	q := New[any]()
	for i := 0; i < b.N; i++ {
		q.Add(nil)
	}
	for i := 0; i < b.N; i++ {
		_, _ = q.Peek()
		_, _ = q.Remove()
	}
}

func BenchmarkQueueGet(b *testing.B) {
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Get(i)
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := New[any]()
	for i := 0; i < b.N; i++ {
		q.Add(nil)
		_, _ = q.Peek()
		_, _ = q.Remove()
	}
}
