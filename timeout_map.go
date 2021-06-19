//Package timeoutmap provides a data structure on top of the
//general map with the addition of providing timeouts for keys
package timeoutmap

import (
	"container/heap"
	"context"
	"fmt"
	"time"
)

//granularity of timeouts is in milliseconds i.e.
//items will be removed withing milliseconds of expiry
type TimeoutMap struct {
	internalMap  map[interface{}]interface{}
	AbsTimeoutPQ durationHeap
	timeout      time.Duration
	ctx          context.Context
}

//Initializes TimeoutMap with a default timeout for all keys
func (tm *TimeoutMap) Init(ctx context.Context, hint int, timeout time.Duration) {
	tm.internalMap = make(map[interface{}]interface{}, hint)
	tm.AbsTimeoutPQ = make(durationHeap, 0)
	heap.Init(&tm.AbsTimeoutPQ)
	tm.timeout = timeout
	tm.ctx = ctx
	go tm.process()
	return
}

//updates the map with the provided key value
//timeout is associated with the key post which the key is removed from the map
//if timeout is 0 then the default timeout is used
//if key already exists then the value is updated and timeout is renewed
func (tm *TimeoutMap) AddKV(key interface{}, value interface{}, timeout time.Duration) {
	tm.internalMap[key] = value
	if timeout == 0 {
		timeout = tm.timeout
	}
	now := time.Now().UnixNano()
	item := &durationItem{
		value:    key,
		priority: int64(tm.timeout) + now,
	}
	heap.Push(&tm.AbsTimeoutPQ, item)
}

//returns value if it exists in the map
//else returns nil
func (tm TimeoutMap) GetValue(key interface{}) (value interface{}, ok bool) {
	if val, ok := tm.internalMap[key]; ok {
		return val, ok
	}
	return nil, false
}

//deletes the key from the map if it exists
func Delete(tm *TimeoutMap, key interface{}) {
	delete(tm.internalMap, key)
}

//returns the length of the map
func Len(tm TimeoutMap) int {
	return len(tm.internalMap)
}

func (tm *TimeoutMap) timeoutHandler() {
	for {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		item := tm.AbsTimeoutPQ.Peek()
		if item == nil {
			break
		}
		if item.(*durationItem).priority/int64(time.Millisecond) < now {
			fmt.Println("deleting")
			heap.Pop(&tm.AbsTimeoutPQ)
			delete(tm.internalMap, item.(*durationItem).value)
		} else {
			break
		}

	}
}

func (tm *TimeoutMap) process() {
	ticker := time.NewTicker(1 * time.Millisecond)
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
