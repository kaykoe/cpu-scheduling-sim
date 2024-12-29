package sim

import (
	"container/ring"
)

type Stack[T any] struct {
	items []T
}

// NewStack returns a pointer to an initialized stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{make([]T, 0, 1024)}
}

func (s *Stack[T]) Push(val T) {
	s.items = append(s.items, val)
}

func (s *Stack[T]) Pop() (val T) {
	val = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return val
}

// Top is a convenience function to get the value of the top element of the stack without the need to pop
func (s *Stack[T]) Top() (val T) {
	return s.items[len(s.items)-1]
}

func (s *Stack[T]) Empty() bool {
	return s.items == nil || len(s.items) == 0
}

type Queue[T any] struct {
	ring *ring.Ring
}

// NewQueue returns a pointer to an initialized queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{ring.New(0)}
}

func (q *Queue[T]) Push(val T) {
	q.ring = q.ring.Link(&ring.Ring{Value: val})
}

func (q *Queue[T]) Pop() (val T) {
	return q.ring.Unlink(1).Value.(T)
}

// Front is a convenience function to get the value of the first element of the queue without the need to pop
func (q *Queue[T]) Front() (val T) {
	return q.ring.Value.(T)
}

func (q *Queue[T]) Empty() bool {
	return q.ring == nil || q.ring.Len() == 0
}
