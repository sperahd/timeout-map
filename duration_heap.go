package timeoutmap

//durationHeap is a min-heap of Durations
//most of the code is take from: https://golang.org/pkg/container/heap/#example__intHeap
//and modified for Duration
type durationItem struct {
	value    interface{}
	priority int64
}

type durationHeap []*durationItem

func (h durationHeap) Len() int {
	return len(h)
}

func (h durationHeap) Less(i, j int) bool {
	return h[i].priority < h[j].priority
}

func (h durationHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *durationHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*durationItem))
}

func (h durationHeap) Peek() interface{} {
	if len(h) <= 0 {
		return nil
	}
	return h[0]
}

func (h *durationHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
