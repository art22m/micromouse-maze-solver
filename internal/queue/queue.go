package queue

import "fmt"

type Queue[T any] struct {
	values []T
}

func (q *Queue[T]) Empty() bool {
	return len(q.values) == 0
}

func (q *Queue[T]) Push(value T) {
	q.values = append(q.values, value)
}

func (q *Queue[T]) Pop() T {
	if len(q.values) == 0 {
		panic("empty queue")
	}

	value := q.values[0]
	q.values = q.values[1:]
	return value
}

func (q *Queue[T]) Print() {
	for _, v := range q.values {
		fmt.Println(v)
	}
}
