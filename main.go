package main

import (
	"fmt"
	"runtime"
	"sync"
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
}
