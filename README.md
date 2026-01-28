# Concurrency in Go

## Goroutines

### Github repo

```sh
git clone https://github.com/andcloudio/go-concurrency-exercises.git
```

### Hello

```go
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
```

### Client / server

```go
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

// handleConn - utility function
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, "response from server\n")
		if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}
```

### WaitGroups

Race condition: execution order not guaranteed

```go
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
```

### Closure 1

```go
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
```

### Closure 2

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait()
}
```

## Go scheduler

### Go scheduler

- M:N scheduler
- User space
- OS threads
- GOMAXPROCS (cpu processors by default)
- Multiple worker threads
- N threads on M processors
- Async preemption
- Max 10ms

#### Goroutines states

- Runnable
- Executing
- Waiting (I/O)

### Context switching

- I/O call
- Thread moved to waiting queue

### Work stealing

Balance work load across processors

## Channels

### Channels

- Communicate data between goroutines
- Synchronize goroutines
- Typed
- Thread-safe
- Pointer operator
- Blocking

### Channel exercise

```go
package main

import "fmt"

func main() {
	ch := make(chan int)

	go func(a, b int) {
		c := a + b
		ch <- c
	}(1, 2)

	// get the value computed from goroutine
	c := <-ch
	fmt.Printf("computed value %v\n", c)
}
```

### Range buffered channles

- Iterates over values received from a channel
- Automatically breaks when channel is closed
- Does not return second value

```go
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
```

### Buffered channels

- given capacity
- in memory FIFO queue
- Async

```go
package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 6)

	go func() {
		defer close(ch)

		// send all iterator values on channel without blocking
		for i := range 6 {
			fmt.Printf("Sending: %d\n", i)
			ch <- i
		}
	}()

	for v := range ch {
		fmt.Printf("Received: %v\n", v)
	}
}
```

### Channel directions

- Only send or receive
- Type-safe

```go
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
```

### Channel ownership

- Default channel value is nil
- Reading/writing from nil is blocking
- Closing nil will panic
- Channel owner is the goroutine that initializes, writes and closes channel
- Channel user is readonly

```go
package main

import "fmt"

func main() {
	// create channel owner goroutine which return channel and
	// - writes data into channel and
	// - closes the channel when done.

	owner := func() <-chan int {
		ch := make(chan int)

		go func() {
			defer close(ch)
			for i := range 6 {
				ch <- i
			}
		}()
		return ch
	}

	consumer := func(ch <-chan int) {
		// read values from channel
		for v := range ch {
			fmt.Printf("Received: %d\n", v)
		}
		fmt.Println("Done receiving!")
	}

	ch := owner()
	consumer(ch)
}
```

## Select

### Select

Select is used to wait on multiple channel operations and run the first one that’s ready.
It’s like a switch, but for channels.

### Select exercise

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "two"
	}()

	// multiplex recv on channel - ch1, ch2
	for range 2 {
		select {
		case m1 := <-ch1:
			fmt.Println(m1)

		case m2 := <-ch2:
			fmt.Println(m2)

		}
	}
}
```

### Select timeout

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "one"
	}()

	// implement timeout for recv on channel ch
	select {
	case m := <-ch:
		fmt.Println(m)
	case <-time.After(3 * time.Second):
		fmt.Println("timeout")
	}

}
```

### Non blocking

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			ch <- "message"
		}

	}()

	// if there is no value on channel, do not block.
	for range 2 {
		select {
		case m := <-ch:
			fmt.Println(m)
		default:
			fmt.Println("No message")
		}

		// Do some processing..
		fmt.Println("processing..")
		time.Sleep(1500 * time.Millisecond)
	}
}
```

## Sync package

### Mutex

- Channels:
  - passing copy of data
  - distributing units of work
  - communicating async results
- Mutex
  - caches
  - states
  - protects shared resources
  - provides exclusive access

### Exercise: mutex

```go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

	runtime.GOMAXPROCS(4)

	var balance int
	var wg sync.WaitGroup
	var mu sync.Mutex

	deposit := func(amount int) {
		mu.Lock()
		defer mu.Unlock()
		balance += amount

	}

	withdrawal := func(amount int) {
		mu.Lock()
		defer mu.Unlock()
		balance -= amount
	}

	// make 100 deposits of $1
	// and 100 withdrawal of $1 concurrently.
	// run the program and check result.

	// TODO: fix the issue for consistent output.

	wg.Add(100)
	for range 100 {
		go func() {
			defer wg.Done()
			deposit(1)
		}()
	}

	wg.Add(100)
	for range 100 {
		go func() {
			defer wg.Done()
			withdrawal(1)
		}()
	}

	wg.Wait()
	fmt.Println(balance)
}
```

### Atomic

- low level operations in memory
- lockless operation
- counters

### Exercise: atomic

```go
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
```

### Conditional variable
