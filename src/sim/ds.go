package sim

import (
	"container/ring"
	"gotest.tools/v3/assert"
	"testing"
)

var t = &testing.T{}

type Stack[T any] struct {
	items []T
}

// NewStack returns a pointer to an initialized stack, only values greater than 0 are allowed for len
func NewStack[T any](len int) *Stack[T] {
	assert.Assert(t, len > 0, "A 0 length stack is nil, initialize with a length")
	return &Stack[T]{make([]T, 0, len)}
}

func (s *Stack[T]) Push(val T) {
	assert.Assert(t, s != nil)
	assert.Assert(t, s.items != nil,
		"The stack's underlying slice should never be nil,"+
			"create the stack with the ds.NewStack() function and not new")

	s.items = append(s.items, val)
}

func (s *Stack[T]) Pop() (val T) {
	assert.Assert(t, s != nil)
	assert.Assert(t, s.items != nil,
		"The stack's underlying slice should never be nil,"+
			"create the stack with the ds.NewStack() function and not new")

	val = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return val
}

// Top is a convenience function to get the value of the top element of the stack without the need to pop
func (s *Stack[T]) Top() (val T) {
	assert.Assert(t, s != nil)
	assert.Assert(t, s.items != nil,
		"The stack's underlying slice should never be nil,"+
			"create the stack with the ds.NewStack() function and not new")

	return s.items[len(s.items)-1]
}

func (s *Stack[T]) Empty() bool {
	assert.Assert(t, s != nil)
	assert.Assert(t, s.items != nil,
		"The stack's underlying slice should never be nil,"+
			"create the stack with the ds.NewStack() function and not new")

	return len(s.items) == 0
}

type Queue[T any] struct {
	ring *ring.Ring
}

// NewQueue returns a pointer to an initialized queue, only values greater than 0 are allowed for len
func NewQueue[T any](len int) *Queue[T] {
	assert.Assert(t, len > 0, "A 0 length Queue is nil, initialize with a length")
	return &Queue[T]{ring.New(len)}
}

func (q *Queue[T]) Push(val T) {
	assert.Assert(t, q != nil)
	assert.Assert(t, q.ring != nil,
		"The queue's underlying ring should never be nil,"+
			"create the queue with the ds.NewQueue() function and not new")
	q.ring = q.ring.Link(&ring.Ring{Value: val})
}

func (q *Queue[T]) Pop() (val T) {
	assert.Assert(t, q != nil)
	assert.Assert(t, q.ring != nil,
		"The queue's underlying ring should never be nil,"+
			"create the queue with the ds.NewQueue() function and not new")

	return q.ring.Unlink(1).Value.(T)
}

// Front is a convenience function to get the value of the first element of the queue without the need to pop
func (q *Queue[T]) Front() (val T) {
	assert.Assert(t, q.ring != nil,
		"The queue's underlying ring should never be nil,"+
			"create the queue with the ds.NewQueue() function and not new")

	return q.ring.Value.(T)
}

func (q *Queue[T]) Empty() bool {
	assert.Assert(t, q.ring != nil,
		"The queue's underlying ring should never be nil,"+
			"create the queue with the ds.NewQueue() function and not new")

	return q.ring.Len() == 0
}
