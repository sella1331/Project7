package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	countPt := 0
	alltasks := 1000
	workers := 5
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	mux := sync.Mutex{}

	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			for {
				mux.Lock()
				if alltasks == 0 {
					mux.Unlock()
					return
				}
				alltasks--
				countPt++
				fmt.Printf("worker #%d: %d\n", id, countPt)
				mux.Unlock()
				runtime.Gosched()
			}
		}(w)
	}

	wg.Wait()
	fmt.Println("count:", countPt)
	test := NewTTLCache()
	test.Set("aлклее", "и шо с того", 5*time.Millisecond)

	if v, ok := test.Get("aлклее"); ok {
		fmt.Println("нашли:", v)
	}
	time.Sleep(5 * time.Second)
	if _, ok := test.Get("aлклее"); !ok {
		fmt.Println("now delete")
		return
	}

}

type timeT struct {
	value any
	died  time.Time
}

type timeLife struct {
	mx sync.RWMutex
	m  map[string]timeT
}

func NewTTLCache() *timeLife {
	return &timeLife{m: make(map[string]timeT)}
}

func (t *timeLife) Get(key string) (any, bool) {
	t.mx.RLock()
	defer t.mx.RUnlock()
	cache, ok := t.m[key]

	if !ok {
		return nil, false
	}

	if time.Now().After(cache.died) {
		return nil, false
	}
	return cache.value, true
}
func (t *timeLife) Set(key string, value any, ttl time.Duration) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.m[key] = timeT{
		value: value,
		died:  time.Now().Add(ttl),
	}
}
func (t *timeLife) Delete(key string) {
	t.mx.Lock()
	defer t.mx.Unlock()
	delete(t.m, key)
}
