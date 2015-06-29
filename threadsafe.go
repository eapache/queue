package queue

import (
	"sync"
)

// A threadsafe queue is Queue structure, wrapped with modifications
// to make it thread safe.
type ThreadSafeQueue struct {
	q    *Queue
	lock *sync.Mutex
}

// Creates and returns a new thread safe queue.
func NewThreadSafe() *ThreadSafeQueue {
	return &ThreadSafeQueue{
		q:    New(),
		lock: new(sync.Mutex),
	}
}

// Length returns the number of elements currently stored in the queue.
func (t *ThreadSafeQueue) Length() int {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.q.Length()
}

// Add puts an element on the end of the queue.
func (t *ThreadSafeQueue) Add(elem interface{}) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.q.Add(elem)
}

// Peek returns the element at the head of the queue. This call errors
// if the queue is empty.
func (t *ThreadSafeQueue) Peek() (interface{}, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.q.Peek()
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will error.
func (t *ThreadSafeQueue) Get(i int) (interface{}, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.q.Get(i)
}

// Gets and returns the first item from the queue.
func (t *ThreadSafeQueue) Pop() (interface{}, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.q.Pop()
}

// Remove removes the element from the front of the queue. If you actually
// want the element, call Peek first. This call errors if the queue is empty.
func (t *ThreadSafeQueue) Remove() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.q.Remove()
}
