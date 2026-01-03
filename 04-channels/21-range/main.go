package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		for id := range 6 {
			// send iterator over channel
			ch <- id
		}
		close(ch)
	}()

	// range over channel to recv values
	for v := range ch {
		fmt.Printf("Val from channel: %d\n", v)
	}

}
