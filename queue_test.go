package queue

import (
	"container/list"
	"testing"
)

const (
	Add    = 0
	Remove = 1
)

type Op struct {
	op    int
	times int
}

var ops = []Op{
	{Add, 50000},
	{Remove, 25000},
	{Add, 200000},
	{Remove, 225000},
	{Add, 10000},
	{Remove, 5000},
	{Add, 1000000},
	{Remove, 1005000},
}

func BenchmarkSliceQueue(b *testing.B) {
	queueBench(&SliceQueue[int]{})
}

func BenchmarkCircularQueue(b *testing.B) {
	queueBench(&CircularQueue[int]{})
}

func queueBench(queue Queue[int]) {
	for _, p := range ops {
		for i := 0; i < p.times; i++ {
			if p.op == Remove {
				_ = queue.Remove()
			} else if p.op == Add {
				queue.Add(0)
			}
		}
	}
}

type Queue[T any] interface {
	Add(T)
	Remove() T
	Empty() bool
}

type SliceQueue[T any] struct {
	q []T
}

func (q *SliceQueue[T]) Add(t T) {
	q.q = append(q.q, t)
}

func (q *SliceQueue[T]) Remove() T {
	t := q.q[0]
	q.q = q.q[1:]
	return t
}

func (q *SliceQueue[T]) Empty() bool {
	return len(q.q) == 0
}

type ListQueue[T any] struct {
	q list.List
}

func (q *ListQueue[T]) Add(t T) {
	q.q.PushBack(t)
}

func (q *ListQueue[T]) Remove() T {
	return q.q.Remove(q.q.Front()).(T)
}

func (q *ListQueue[T]) Empty() bool {
	return q.q.Len() == 0
}

type CircularQueue[T any] struct {
	q     []T
	front int
	back  int
	size  int
}

func (q *CircularQueue[T]) Add(t T) {
	if q.size == len(q.q) {
		q.resize()
	}
	q.q[q.back] = t
	q.back = (q.back + 1) % len(q.q)
	q.size++
}

func (q *CircularQueue[T]) Remove() T {
	t := q.q[q.front]
	q.front = (q.front + 1) % len(q.q)
	q.size--
	return t
}

func (q *CircularQueue[T]) Empty() bool {
	return q.size == 0
}

func (q *CircularQueue[T]) Clear() {
	q.q = nil
	q.front = 0
	q.back = 0
	q.size = 0
}

func (q *CircularQueue[T]) resize() {
	var newCap int
	if q.q == nil {
		newCap = 1
	} else if q.size <= 1024 {
		newCap = len(q.q) << 1
	} else {
		newCap = len(q.q)
		newCap += len(q.q) >> 1
		newCap += len(q.q) >> 3
		newCap += len(q.q) >> 5
	}
	s := make([]T, newCap)
	m := copy(s, q.q[q.front:len(q.q)])
	n := copy(s[m:], q.q[0:q.front])
	if m+n != q.size {
		panic("incorrect size")
	}
	q.q = s
	q.front = 0
	q.back = q.size
}
