package stack

import (
	"sync"
)

type Stack struct {
	top   *element
	size  int
	mutex *sync.Mutex
}

type element struct {
	value interface{}
	next  *element
}

// Return the stack's length
func (s *Stack) Len() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.size
}

// Push a new element onto the stack
func (s *Stack) Push(value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.top = &element{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *Stack) Pop() (value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

func NewStack() *Stack {
	s := new(Stack)
	s.mutex = new(sync.Mutex)
	s.size = 0
	s.top = nil
	return s
}
