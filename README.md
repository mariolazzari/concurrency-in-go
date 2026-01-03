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
