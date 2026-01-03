package main

import (
	"fmt"
	"sync"
)

func main() {
	var data int
	var wg sync.WaitGroup

	wg.Go(func() {
		data++
	})

	wg.Wait()
	fmt.Printf("Data is %v\n", data)
	fmt.Println("Done")
}
