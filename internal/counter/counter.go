package counter

import (
	"math"

	"github.com/emirpasic/gods/trees/binaryheap"
)

type Counter[T uint | float32] struct {
	counter map[string]T
}

type counterEntry[T uint | float32] struct {
	Key   string
	Count T
}

func New[T uint | float32]() *Counter[T] {
	return &Counter[T]{
		counter: map[string]T{},
	}
}

func (c *Counter[T]) Top(n int) []counterEntry[T] {
	heap := *binaryheap.NewWith(c.compare)
	for k, v := range c.counter {
		heap.Push(counterEntry[T]{k, v})
	}
	entries := []counterEntry[T]{}
	for {
		v, ok := heap.Pop()
		if !ok {
			break
		}
		entry := v.(counterEntry[T])
		entries = append(entries, entry)
		if len(entries) >= n {
			return entries
		}
	}
	return entries
}

func (c *Counter[T]) Average() float64 {
	total := len(c.counter)
	if total == 0 {
		return math.NaN()
	}
	sum := T(0)
	for _, v := range c.counter {
		sum += v
	}
	return float64(sum) / float64(total)
}

func (c *Counter[T]) compare(a, b interface{}) int {
	x := a.(counterEntry[T]).Count
	y := b.(counterEntry[T]).Count
	if y > x {
		return 1
	}
	if x < y {
		return -1
	}
	return 0
}

func (c *Counter[T]) Increment(key string, count T) {
	v, ok := c.counter[key]
	if !ok {
		v = 0
	}
	value := v + count
	c.counter[key] = value
}

func (c *Counter[T]) Count(key string) T {
	v, ok := c.counter[key]
	if !ok {
		return 0
	}
	return v
}
