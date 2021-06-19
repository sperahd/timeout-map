package timeoutmap

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func localCheck(got, want interface{}) bool {
	switch got.(type) {
	case string:
		if got != want {
			return false
		}
	}
	return true
}

func TestDefaultTimeout(t *testing.T) {
	defer goleak.VerifyNone(t)
	tm := &TimeoutMap{}
	ctx, cancelFunc := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)
	tm.Init(ctx, 0, time.Duration(2*time.Second), wg)
	key := "id1"
	value := "value1"
	tm.Store(key, value, 0)

	want := value
	got, _ := tm.Load(key)
	if !localCheck(got, want) {
		panic(fmt.Errorf("failed, got: %v, want: %v", got, want))
	}
	time.Sleep(3 * time.Second)

	// timeout case, value should not exist
	if _, ok := tm.Load("id1"); ok {
		panic("failed, item not removed even after timeout")
	}

	cancelFunc()
	wg.Wait()
}
