package sim

import (
	"container/ring"
	"log"
)

type Stack[T any] struct {
	items []T
}

// NewStack returns a pointer to an initialized stack, only values greater than 0 are allowed for len
func NewStack[T any](len int) *Stack[T] {
	if len <= 0 {
		log.Panic("A 0 length stack is nil, initialize with a length")
	}
	return &Stack[T]{make([]T, 0, len)}
}

func (s *Stack[T]) Push(val T) {
	if s == nil {
		log.Panic("pointer to stack cannot be nil")
	}
	if s.items == nil {
		log.Panic("The stack's underlying slice should never be nil,",
			"\ncreate the stack with the ds.NewStack() function and not new")
	}

	s.items = append(s.items, val)
}

func (s *Stack[T]) Pop() (val T) {
	if s == nil {
		log.Panic("pointer to stack cannot be nil")
	}
	if s.items == nil {
		log.Panic("The stack's underlying slice should never be nil,",
			"\ncreate the stack with the ds.NewStack() function and not new")
	}

	val = s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return val
}

// Top is a convenience function to get the value of the top element of the stack without the need to pop
func (s *Stack[T]) Top() (val T) {
	if s == nil {
		log.Panic("pointer to stack cannot be nil")
	}
	if s.items == nil {
		log.Panic("The stack's underlying slice should never be nil,",
			"\ncreate the stack with the ds.NewStack() function and not new")
	}

	return s.items[len(s.items)-1]
}

func (s *Stack[T]) Empty() bool {
	if s == nil {
		log.Panic("pointer to stack cannot be nil")
	}
	if s.items == nil {
		log.Panic("The stack's underlying slice should never be nil,",
			"\ncreate the stack with the ds.NewStack() function and not new")
	}

	return len(s.items) == 0
}

type Queue[T any] struct {
	ring *ring.Ring
}

// NewQueue returns a pointer to an initialized queue, only values greater than 0 are allowed for len
func NewQueue[T any](len int) *Queue[T] {
	if len <= 0 {
		log.Panic("A 0 length Queue is nil, initialize with a length")
	}

	return &Queue[T]{ring.New(len)}
}

func (q *Queue[T]) Push(val T) {
	if q == nil {
		log.Panic("pointer to queue cannot be nil")
	}
	if q.ring == nil {
		log.Panic("The queue's underlying ring should never be nil,",
			"\ncreate the queue with the ds.NewQueue() function and not new")
	}

	q.ring = q.ring.Link(&ring.Ring{Value: val})
}

func (q *Queue[T]) Pop() (val T) {
	if q == nil {
		log.Panic("pointer to queue cannot be nil")
	}
	if q.ring == nil {
		log.Panic("The queue's underlying ring should never be nil,",
			"\ncreate the queue with the ds.NewQueue() function and not new")
	}

	return q.ring.Unlink(1).Value.(T)
}

// Front is a convenience function to get the value of the first element of the queue without the need to pop
func (q *Queue[T]) Front() (val T) {
	if q == nil {
		log.Panic("pointer to queue cannot be nil")
	}
	if q.ring == nil {
		log.Panic("The queue's underlying ring should never be nil,",
			"\ncreate the queue with the ds.NewQueue() function and not new")
	}

	return q.ring.Value.(T)
}

func (q *Queue[T]) Empty() bool {
	if q == nil {
		log.Panic("pointer to queue cannot be nil")
	}
	if q.ring == nil {
		log.Panic("The queue's underlying ring should never be nil,",
			"\ncreate the queue with the ds.NewQueue() function and not new")
	}

	return q.ring.Len() == 0
}
