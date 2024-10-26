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

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.values) == 0 {
		var t T
		return t, false
	}
	value := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return value, true
}

func (s *Stack[T]) Print() {
	for i := len(s.values) - 1; i >= 0; i-- {
		fmt.Println(s.values[i])
	}
}
