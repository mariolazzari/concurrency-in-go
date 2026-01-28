package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

func main() {
	runtime.GOMAXPROCS(4)

	var counter uint64
	var wg sync.WaitGroup

	// implement concurrency safe counter

	for range 50 {
		wg.Go(func() {
			for range 1000 {
				// counter++
				atomic.AddUint64(&counter, 1)
			}
		})
	}
	wg.Wait()
	fmt.Println("counter: ", counter)
}
