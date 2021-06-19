package timeoutmap

/*
package timeoutmap // import "github.com/sperahd/timeout-map"

Package timeoutmap provides a data structure on top of the general map with
the addition of providing timeouts for keys

FUNCTIONS

func Delete(tm *TimeoutMap, key interface{})
    deletes the key from the map if it exists

func Len(tm TimeoutMap) int
    returns the length of the map


TYPES

type TimeoutMap struct {
	AbsTimeoutPQ durationHeap
	// Has unexported fields.
}
    granularity of timeouts is in milliseconds i.e. items will be removed
    withing milliseconds of expiry

func (tm *TimeoutMap) AddKV(key interface{}, value interface{}, timeout time.Duration)
    updates the map with the provided key value timeout is associated with the
    key post which the key is removed from the map if timeout is 0 then the
    default timeout is used if key already exists then the value is updated and
    timeout is renewed

func (tm TimeoutMap) GetValue(key interface{}) (value interface{}, ok bool)
    returns value if it exists in the map else returns nil

func (tm *TimeoutMap) Init(hint int, timeout time.Duration)
    Initializes TimeoutMap with a default timeout for all keys
*/
