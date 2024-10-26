package stack

import (
	"fmt"
)

type Stack[T any] struct {
	values []T
}

func (s *Stack[T]) Empty() bool {
	return len(s.values) == 0
}

func (s *Stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *Stack[T]) Pop() T {
	if len(s.values) == 0 {
		panic("empty stack")
	}
	value := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return value
}

func (s *Stack[T]) Print() {
	for i := len(s.values) - 1; i >= 0; i-- {
		fmt.Println(s.values[i])
	}
}
