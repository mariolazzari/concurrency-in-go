package main

import "fmt"

// Implement relaying of message with Channel Direction

func genMsg(ch chan<- string) {
	defer close(ch)
	// send message on ch1
	ch <- "Ciao Mario"
}

func relayMsg(in <-chan string, out chan<- string) {
	defer close(out)
	// recv message on ch1
	msg := <-in
	// send it on ch2
	out <- msg
}

func main() {
	// create ch1 and ch2
	ch1 := make(chan string)
	ch2 := make(chan string)

	// spine goroutines genMsg and relayMsg
	go genMsg(ch1)
	go relayMsg(ch1, ch2)

	// recv message on ch2
	v := <-ch2
	fmt.Println(v)
}
