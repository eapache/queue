package queue

import "testing"

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
