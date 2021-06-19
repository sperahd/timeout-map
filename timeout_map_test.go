package timeoutmap

import (
	"context"
	"fmt"
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
	tm.Init(ctx, 0, time.Duration(2*time.Second))
	key := "id1"
	value := "value1"
	tm.AddKV(key, value, 0)

	want := value
	got, _ := tm.GetValue(key)
	if !localCheck(got, want) {
		panic(fmt.Errorf("failed, got: %v, want: %v", got, want))
	}
	time.Sleep(3 * time.Second)

	// timeout case, value should not exist
	if _, ok := tm.GetValue("id1"); ok {
		panic("failed, item not removed even after timeout")
	}

	if Len(*tm) != 0 {
		panic("length should be zero")
	}
	cancelFunc()
	time.Sleep(5 * time.Second)
}
