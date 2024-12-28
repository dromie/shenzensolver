package util

import (
	"container/heap"
)

type Item[T any] struct {
	Value    T
	Priority int
}

type internalItem[T any] struct {
	item  *Item[T]
	index int
}

type internalPriorityQueue[T any] []*internalItem[T]

func (pq internalPriorityQueue[T]) Len() int { return len(pq) }

func (pq internalPriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].item.Priority > pq[j].item.Priority
}

func (pq internalPriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *internalPriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*internalItem[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *internalPriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *internalPriorityQueue[T]) update(item *internalItem[T], value T, priority int) {
	item.item.Value = value
	item.item.Priority = priority
	heap.Fix(pq, item.index)
}

type PriorityQueue[T any] struct {
	internal internalPriorityQueue[T]
}

func Pqueue_init[T any]() PriorityQueue[T] {
	pq := PriorityQueue[T]{}
	pq.internal = internalPriorityQueue[T]{}
	heap.Init(&pq.internal)
	return pq
}

func (pq *PriorityQueue[T]) Len() int {
	return pq.internal.Len()
}

func (pq *PriorityQueue[T]) Push(x *Item[T]) {
	heap.Push(&pq.internal, &internalItem[T]{item: x, index: pq.Len()})
}

func (pq *PriorityQueue[T]) Pop() *Item[T] {
	return heap.Pop(&pq.internal).(*internalItem[T]).item
}
