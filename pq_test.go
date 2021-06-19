package timeoutmap

import (
	"container/heap"
	"fmt"
	"testing"
)

func check(got, want interface{}) {
	if got.(*HeapItem) != want.(*HeapItem) {
		fmt.Printf("got: %v, want: %v\n", got, want)
		panic("failed")
	}
}

func TestPQ(t *testing.T) {
	durationPQ := make(PQ, 0)
	heap.Init(&durationPQ)

	// Push-Pop test cases
	heapItem1 := &HeapItem{
		value:    "item1",
		priority: 15,
	}

	heapItem2 := &HeapItem{
		value:    "item2",
		priority: 20,
	}

	heapItem3 := &HeapItem{
		value:    "item3",
		priority: 10,
	}

	heap.Push(&durationPQ, heapItem1)
	heap.Push(&durationPQ, heapItem2)

	want := heapItem1
	got := heap.Pop(&durationPQ)
	check(got, want)

	// Push-Peek-Pop test cases
	heap.Push(&durationPQ, heapItem3)

	want = heapItem3
	got = durationPQ.Peek()
	check(got, want)

	got = heap.Pop(&durationPQ)
	check(got, want)

	want = heapItem2
	got = heap.Pop(&durationPQ)
	check(got, want)

	// Push-Remove-Pop test cases
	heap.Push(&durationPQ, heapItem1)
	heap.Push(&durationPQ, heapItem2)
	heap.Push(&durationPQ, heapItem3)

	want = heapItem2
	got = durationPQ.RemoveItem(heapItem2)
	check(got, want)

	want = heapItem3
	got = heap.Pop(&durationPQ)
	check(got, want)

	want = heapItem1
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
