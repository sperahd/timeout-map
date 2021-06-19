package timeoutmap

import (
	"container/heap"
	"fmt"
	"testing"
)

func check(got, want interface{}) {
	if got.(*durationItem) != want.(*durationItem) {
		fmt.Printf("got: %v, want: %v\n", got, want)
		panic("failed")
	}
}

func TestDurationHeap(t *testing.T) {
	durationPQ := make(durationHeap, 0)
	heap.Init(&durationPQ)

	heapItem1 := &durationItem{
		value:    "item1",
		priority: 15,
	}

	heapItem2 := &durationItem{
		value:    "item2",
		priority: 20,
	}

	heapItem3 := &durationItem{
		value:    "item3",
		priority: 10,
	}

	heap.Push(&durationPQ, heapItem1)
	heap.Push(&durationPQ, heapItem2)

	want := heapItem1
	got := heap.Pop(&durationPQ)
	check(got, want)

	heap.Push(&durationPQ, heapItem3)

	want = heapItem3
	got = durationPQ.Peek()
	check(got, want)

	got = heap.Pop(&durationPQ)
	check(got, want)

	want = heapItem2
	got = heap.Pop(&durationPQ)
	check(got, want)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered")
		} else {
			panic("panic missing for out of range access")
		}
	}()

	heap.Pop(&durationPQ)
}
