package timeoutmap

import (
	"container/heap"
)

//PQ is a min-heap prioritized on int64 values
//reference code is take from: https://golang.org/pkg/container/heap/#example__intHeap
//and modified further to add more APIs
type HeapItem struct {
	value    interface{}
	priority int64
	index    int
}

type PQ []*HeapItem

func (h PQ) Len() int {
	return len(h)
}

func (h PQ) Less(i, j int) bool {
	return h[i].priority < h[j].priority
}

func (h PQ) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *PQ) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	n := len(*h)
	item := x.(*HeapItem)
	item.index = n
	*h = append(*h, item)
}

func (h *PQ) RemoveItem(x interface{}) interface{} {
	return h.Remove(x.(*HeapItem).index)

}

func (h *PQ) Remove(index int) interface{} {
	if index >= len(*h) {
		return nil
	}
	return heap.Remove(h, index)
}

func (h PQ) Peek() interface{} {
	if len(h) <= 0 {
		return nil
	}
	return h[0]
}

func (h *PQ) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
