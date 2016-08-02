/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package queue

// minQueueLen is smallest capacity that queue may have.
// Must be power of 2 for modular arithmetic.
const minQueueLen = 16

// Queue represents a single instance of the queue data structure.
type Queue struct {
	buf               []interface{}
	head, tail, count int
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

// Capacity returns the maximum number of elements that the queue can store
// before increasing allocated storage.
func (q *Queue) Capacity() int {
	return len(q.buf)
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
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	// Compute modulus bitwise since len(q.buf) is always a power of 2:
	// x % n == x & (n - 1)
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Queue) Peek() interface{} {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.buf[q.head]
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Queue) Get(i int) interface{} {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i = q.count + i
	}
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Remove removes the element from the front of the queue. If you actually
// want the element, call Peek first. This call panics if the queue is empty.
func (q *Queue) Remove() {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	q.buf[q.head] = nil
	// Compute modulus bitwise since len(q.buf) is always a power of 2:
	// x % n == x & (n - 1)
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
}

// Pop removes and returns the element from the front of (first in) the queue.
// If the queue is empty, returns nil and false to indicate the queue was empty
// and no data was returned.
func (q *Queue) Pop() (interface{}, bool) {
	if q.count <= 0 {
		return nil, false
	}
	ret := q.buf[q.head]
	q.buf[q.head] = nil
	// Compute modulus bitwise since len(q.buf) is always a power of 2:
	// x % n == x & (n - 1)
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret, true
}

// Clear removes all elements from the queue, but retains the current capacity.
func (q *Queue) Clear() {
	// Compute modulus bitwise since len(q.buf) is always a power of 2:
	// x % n == x & (n - 1)
	blen := len(q.buf) - 1
	for h := q.head; h != q.tail; h = (h + 1) & blen {
		q.buf[h] = nil
	}
	q.head = 0
	q.tail = 0
	q.count = 0
}
