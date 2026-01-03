package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	incr := func(wg *sync.WaitGroup) {
		var i int
		wg.Go(func() {
			i++
			fmt.Printf("value of i: %v\n", i)
		})
		fmt.Println("return from function")
	}

	incr(&wg)
	wg.Wait()
	fmt.Println("done..")
}
