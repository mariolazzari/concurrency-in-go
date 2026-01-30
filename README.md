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

- sync mechanism
- goroutines container
- wait suspends execution of a goroutine
- signal wakes one goroutine
- broadcast wakes up all go routines

### Exercise: signal

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedRsc = make(map[string]any)

func main() {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	wg.Go(func() {

		// suspend goroutine until sharedRsc is populated.
		c.L.Lock()

		for len(sharedRsc) == 0 {
			time.Sleep(1 * time.Millisecond)
		}

		fmt.Println(sharedRsc["rsc1"])
		c.L.Unlock()
	})

	c.L.Lock()
	// writes changes to sharedRsc
	sharedRsc["rsc1"] = "foo"
	c.Signal()
	c.L.Unlock()

	wg.Wait()
}
```

### Exercise: broadcast

```go
package main

import (
	"fmt"
	"sync"
)

var sharedRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	wg.Go(func() {

		// suspend goroutine until sharedRsc is populated.
		c.L.Lock()

		for len(sharedRsc) == 0 {
			c.L.Lock()
		}

		fmt.Println(sharedRsc["rsc1"])
		c.L.Unlock()
	})

	wg.Add(1)
	go func() {
		defer wg.Done()

		// suspend goroutine until sharedRsc is populated.
		c.L.Lock()

		for len(sharedRsc) == 0 {
			c.Wait()
		}

		fmt.Println(sharedRsc["rsc2"])
		c.L.Unlock()
	}()

	c.L.Lock()
	// writes changes to sharedRsc
	sharedRsc["rsc1"] = "foo"
	sharedRsc["rsc2"] = "bar"
	c.Broadcast()
	c.L.Unlock()

	wg.Wait()
}
```

### Sync once

- one time init function

```go
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
```

### Sync pool

- pool for expensive resources

### Exercise: pool

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// create pool of bytes.Buffers which can be reused.
var bufPool = sync.Pool{
	New: func() any {
		fmt.Println("Allocating new buffer")
		return new(bytes.Buffer)
	},
}

func log(w io.Writer, val string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()

	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(val)
	b.WriteString("\n")

	w.Write(b.Bytes())
	bufPool.Put(b)
}

func main() {
	log(os.Stdout, "debug-string1")
	log(os.Stdout, "debug-string2")
}
```

## Race detector

- tool for finding race conditions
- 10 times slower
- use in test only

```sh
go run -race main.go
```

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	var t *time.Timer

	ch := make(chan bool)

	t = time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Since(start))
		ch <- true
	})

	for time.Since(start) < 5*time.Second {
		<-ch
		t.Reset(randomDuration())

	}

	time.Sleep(5 * time.Second)
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}
```

### Web crawler

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

var fetched map[string]bool

type result struct {
	url   string
	urls  []string
	err   error
	depth int
}

// Crawl uses findLinks to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int) {
	results := make(chan *result)

	fetch := func(url string, depth int) {
		urls, err := findLinks(url)
		results <- &result{url, urls, err, depth}
	}

	go fetch(url, depth)
	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		res := <-results
		if res.err != nil {
			// fmt.Println(res.err)
			continue
		}

		fmt.Printf("found: %s\n", res.url)
		if res.depth > 0 {
			for _, u := range res.urls {
				if !fetched[u] {
					fetching++
					go fetch(u, res.depth-1)
					fetched[u] = true
				}
			}
		}
	}
	close(results)
}

func main() {
	fetched = make(map[string]bool)
	now := time.Now()
	Crawl("http://andcloud.io", 2)
	fmt.Println("time taken:", time.Since(now))
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

// visit appends to links each link found in n, and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
```

## Concurrency patterns

### Pipelines

- process streams
- more stages
- concurrent stages
- same type
- composable

### Exercise: pipelines

```go
package main

import "fmt"

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// set up the pipeline
	for n := range square(square(generator(2, 3))) {
		fmt.Println(n)
	}
}
```

### Fan out, fan in

- Brakes cpu intensive stage into more goroutines
- fan out: multiple goroutines read from single chanel
- fan in: comining multiple results to a signle channel

### Exercise: fan out, fan in

```go
// Squaring numbers.

package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	in := generator(2, 3)

	c1 := square(in)
	c2 := square(in)

	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}
```

### Cancelling goroutines

- upstreams close channels
- downstreams receave until channel is close
- readonly done channel to close

### Exercise: cancelling goroutines

```go
// Squaring numbers.

package main

import (
	"fmt"
	"sync"
)

func generator(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}

	}()
	return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	in := generator(done, 2, 3)

	c1 := square(done, in)
	c2 := square(done, in)

	out := merge(done, c1, c2)

	fmt.Println(<-out)
}
```

### Image processing pipeline

```go
package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

// Pipeline
// walkfile ----------> processImage -----------> saveImage
//            (paths)                  (results)

type result struct {
	srcImagePath   string
	thumbnailImage *image.NRGBA
	err            error
}

// Image processing - Pipeline
// Input - directory with images.
// output - thumbnail images
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}
	start := time.Now()
	err := setupPipeLine(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func setupPipeLine(root string) error {
	done := make(chan struct{})
	defer close(done)

	// do the file walk
	paths, errc := walkFiles(done, root)

	// process the image
	results := processImage(done, paths)

	// save thumbnail images
	for r := range results {
		if r.err != nil {
			return r.err
		}
		saveThumbnail(r.srcImagePath, r.thumbnailImage)
	}

	// check for error on the channel, from walkfiles stage.
	if err := <-errc; err != nil {
		return err
	}
	return nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {

	// create output channels
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			// filter out error
			if err != nil {
				return err
			}

			// check if it is file
			if !info.Mode().IsRegular() {
				return nil
			}

			// check if it is image/jpeg
			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return nil
			}

			// send file path to next stage
			select {
			case paths <- path:
			case <-done:
				return fmt.Errorf("walk cancelled")
			}
			return nil
		})
	}()
	return paths, errc
}

func processImage(done <-chan struct{}, paths <-chan string) <-chan result {
	results := make(chan result)
	var wg sync.WaitGroup

	thumbnailer := func() {
		for srcImagePath := range paths {
			srcImage, err := imaging.Open(srcImagePath)
			if err != nil {
				select {
				case results <- result{srcImagePath, nil, err}:
				case <-done:
					return
				}
			}
			thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

			select {
			case results <- result{srcImagePath, thumbnailImage, err}:
			case <-done:
				return
			}
		}
	}

	const numThumbnailer = 5
	for i := 0; i < numThumbnailer; i++ {
		wg.Add(1)
		go func() {
			thumbnailer()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

// saveThumbnail - save the thumnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

// getFileContentType - return content type and error status
func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
```

## Context package

### Context

- propagate request scope
- propagate cancellation signal

#### Background

- empty context
- root of context tree
- never cancelled
- no values
- no deadline
- used by main function
- top level of incoming request

#### Todo

- root of context tree
- placeholder
-

### Cancellation

### Data bag
