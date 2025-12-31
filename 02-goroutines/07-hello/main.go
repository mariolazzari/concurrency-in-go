package main

import (
	"fmt"
	"time"
)

func fun(s string) {
	for range 3 {
		fmt.Println(s)
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	// Direct call
	fun("direct call")

	// goroutine function call
	go fun("goroutine-1")

	// goroutine with anonymous function
	go func() {
		fun("goroutine-2")
	}()

	// goroutine with function value call
	fv := fun
	go fv("goroutine-3")

	// wait for goroutines to end
	fmt.Println("Wait for goroutines")
	time.Sleep(100 * time.Millisecond)

	fmt.Println("done..")
}
