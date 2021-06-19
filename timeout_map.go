//Package timeoutmap provides a data structure on top of the
//general map with the addition of providing timeouts for keys
package timeoutmap

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

//granularity of timeouts is in milliseconds i.e.
//items will be removed withing milliseconds of expiry
type TimeoutMap struct {
	internalMap  *sync.Map
	AbsTimeoutPQ PQ
	timeout      time.Duration
	ctx          context.Context
	pqMutex      *sync.Mutex
	wg           *sync.WaitGroup
}

//Initializes TimeoutMap with a default timeout for all keys
func (tm *TimeoutMap) Init(ctx context.Context, hint int, timeout time.Duration, wg *sync.WaitGroup) {
	tm.internalMap = new(sync.Map)
	tm.AbsTimeoutPQ = make(PQ, 0)
	heap.Init(&tm.AbsTimeoutPQ)
	tm.timeout = timeout
	tm.ctx = ctx
	tm.pqMutex = new(sync.Mutex)
	tm.wg = wg
	tm.wg.Add(1)
	go tm.process()
	return
}

//updates the map with the provided key value
//timeout is associated with the key post which the key is removed from the map
//if timeout is 0 then the default timeout is used
//if key already exists then the value is updated and timeout is renewed
func (tm *TimeoutMap) Store(key interface{}, value interface{}, timeout time.Duration) {

	// TODO: Update this to update the timeout of the existing
	// item if it exists
	// Will need update priority API in the pq implementation
	tm.internalMap.Store(key, value)
	if timeout == 0 {
		timeout = tm.timeout
	}
	now := time.Now().UnixNano()
	item := &HeapItem{
		value:    key,
		priority: int64(tm.timeout) + now,
	}
	tm.pqMutex.Lock()
	defer tm.pqMutex.Unlock()
	heap.Push(&tm.AbsTimeoutPQ, item)
}

//returns value if it exists in the map
//else returns nil
func (tm *TimeoutMap) Load(key interface{}) (value interface{}, ok bool) {
	return tm.internalMap.Load(key)
}

//deletes the key from the map if it exists
func (tm *TimeoutMap) Delete(key interface{}) {
	tm.internalMap.Delete(key)
	tm.pqMutex.Lock()
	defer tm.pqMutex.Unlock()
	tm.AbsTimeoutPQ.RemoveItem(key)
}

func (tm *TimeoutMap) timeoutHandler() {
	for {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		// check for existence of an item
		tm.pqMutex.Lock()
		item := tm.AbsTimeoutPQ.Peek()
		tm.pqMutex.Unlock()
		if item == nil {
			break
		}

		// check for timeout expiry against current time
		if now > item.(*HeapItem).priority/int64(time.Millisecond) {
			fmt.Println("deleting")
			tm.pqMutex.Lock()
			heap.Pop(&tm.AbsTimeoutPQ)
			tm.internalMap.Delete(item.(*HeapItem).value)
			tm.AbsTimeoutPQ.RemoveItem(item)
			tm.pqMutex.Unlock()
		} else {
			break
		}

	}
}

func (tm *TimeoutMap) process() {
	ticker := time.NewTicker(1 * time.Millisecond)

	// following is executed in lifo order
	defer tm.wg.Done()
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ticker.C:
			tm.timeoutHandler()
		case <-tm.ctx.Done():
			break loop
		}
	}
}
