package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var once sync.Once

	load := func() {
		fmt.Println("Run only once initialization function")
	}

	for range 10 {
		go func() {
			defer wg.Done()

			// modify so that load function gets called only once.
			once.Do(load)
		}()
	}
	wg.Wait()
}
