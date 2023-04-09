/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package queue

import (
	"errors"
	"sync"
)

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
const minQueueLen = 16

var (
	ErrQueueEmpty      = errors.New("queue is empty")
	ErrIndexOutOfRange = errors.New("index is out of range")
)

// Queue represents a single instance of the queue data structure.
type Queue struct {
	buf               []interface{}
	head, tail, count int
	mu                sync.Mutex
}

// New constructs and returns a new Queue.
func New() *Queue {
	return &Queue{
		buf: make([]interface{}, minQueueLen),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Queue) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newBuf := make([]interface{}, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *Queue) Add(elem interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This return error
// if the queue is empty.
func (q *Queue) Peek() (interface{}, error) {
	if q.count <= 0 {
		return nil, ErrQueueEmpty
	}
	return q.buf[q.head], nil
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will return error. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue) Get(i int) (interface{}, error) {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		return nil, ErrIndexOutOfRange
	}
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)], nil
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will return error.
func (q *Queue) Remove() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.count <= 0 {
		return nil, ErrQueueEmpty
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.

	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret, nil
}
