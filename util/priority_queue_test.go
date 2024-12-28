package util

import (
	"testing"
)

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func Test_pqueue(t *testing.T) {
	// Some test_items and their priorities.
	test_items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := Pqueue_init[string]()
	for value, priority := range test_items {
		pq.Push(&Item[string]{
			Value:    value,
			Priority: priority,
		})
	}
	if pq.Len() != 3 {
		t.Errorf("Priority queue length is not 3")
	}
	// Insert a new item and then modify its priority.
	item := &Item[string]{
		Value:    "orange",
		Priority: 5,
	}
	pq.Push(item)
	if pq.Len() != 4 {
		t.Errorf("Priority queue length is not 4")
	}

	firstItem := pq.Pop()
	if firstItem.Value != "orange" {
		t.Errorf("First item is not orange")
	}
	secondItem := pq.Pop()
	if secondItem.Value != "pear" {
		t.Errorf("Second item is not pear")
	}
	thirdItem := pq.Pop()
	if thirdItem.Value != "banana" {
		t.Errorf("Third item is not banana")
	}
	fourthItem := pq.Pop()
	if fourthItem.Value != "apple" {
		t.Errorf("Fourth item is not apple")
	}
	if pq.Len() != 0 {
		t.Errorf("Priority queue length is not 0")
	}
}
