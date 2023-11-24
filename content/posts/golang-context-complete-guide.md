---
layout: post
title: Golang Context Guide in Concurrent Programs 
summary: Let's check how to use context effectively in concurrency. 
date: 2023-11-24
tags: [golang, concurrency]
---

Concurrent programming in Go can be a bit like juggling—keeping many tasks in the air at once. The `context` package in Go acts as your trusty assistant, helping you manage this juggling act with finesse.

In this blog post, we're going to unravel the secrets of the `context` package. It's your toolkit for handling tricky situations in concurrent programs, such as canceling tasks, setting deadlines, and smoothly passing information between different parts of your code.

Think of it as your guide to becoming a concurrency maestro in Go. We'll start with the basics, explore how to use the context package in real-life scenarios, and wrap up with tips to keep your concurrent programs running smoothly.

## Context Package

Let's have a look at the code for the `context` package.

```go
// A Context carries a deadline, a cancellation signal, and other values across
// API boundaries.
//
// Context's methods may be called by multiple goroutines simultaneously.
type Context interface {
	// Deadline returns the time when work done on behalf of this context
	// should be canceled. Deadline returns ok==false when no deadline is
	// set. Successive calls to Deadline return the same results.
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel that's closed when work done on behalf of this
	// context should be canceled. Done may return nil if this context can
	// never be canceled. Successive calls to Done return the same value.
	// The close of the Done channel may happen asynchronously,
	// after the cancel function returns.
	Done() <-chan struct{}

	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error explaining why:
	// Canceled if the context was canceled
	// or DeadlineExceeded if the context's deadline passed.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	Err() error

    // Value returns the value associated with this context for key, or nil
	// if no value is associated with key. Successive calls to Value with
	// the same key returns the same result.
	Value(key any) any
}
```

It is quite simple and we will be mostly using the `Done() <- chan struct{}` function.

Let's have a look at the functions provided by the `context` package that we will use.

```go
// WithCancel returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called
// or when the parent context's Done channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d. If the parent's deadline is already earlier than d,
// WithDeadline(parent, d) is semantically equivalent to parent. The returned
// [Context.Done] channel is closed when the deadline expires, when the returned
// cancel function is called, or when the parent context's Done channel is
// closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this [Context] complete.
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this [Context] complete:
//
//	func slowOperationWithTimeout(ctx context.Context) (Result, error) {
//		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
//		defer cancel()  // releases resources if slowOperation completes before timeout elapses
//		return slowOperation(ctx)
//	}
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

In most of the times these 3 functions are used and I will be giving examples on how to use them effectively.

## Demo Showcase of the Concurrency

We will write some pseudo-code first to understand the flow, then implement it in golang.

```
- channel is declared
- some goroutine(s) is writing into the channel
- some goroutines(s) are listening from channel in parallel and process the messages
- wait for the listeners to finish the processing and exit the program
```

For our example, we will be having 2 goroutines listening to a channel and the main goroutine writing into the channel.

A simple program is with 2 second timeout can be written as

```go
package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dataChannel := make(chan int)

	var wg sync.WaitGroup

	wg.Add(2)
	for j := 0; j < 2; j++ {
		go func(i int) {
			worker(ctx, dataChannel, i)
			wg.Done()
		}(j)
	}

	i := 0
	for i < 5 {
		dataChannel <- i
		i++
	}

	wg.Wait()
}

func worker(ctx context.Context, ch chan int, number int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("context canceled, exiting for worker number %d\n", number)
			return
		case i, ok := <-ch:
			// It means channel is closed
			if !ok {
				log.Printf("channel is closed, exiting worker %d\n", number)
				return
			}
			log.Printf("worker: %d: read %v from channel\n", number, i)
		}
	}
}
```

### Explanation of the Demo

Let's split the code and explain what each part is doing.

1. 
```go
// ...
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

dataChannel := make(chan int)
// ...
```

We are initializing the context variable with 2 second timeout. So if our program is taking more than 2 second, we want the 
program to exit it and inform the sub-goroutines about it.

Calling `defer cancel()`  releases resources if the program completes before 2 second.

In the last line, we are creating the channel which will be the main communication system for our program.

2. 
```go
// ...
	var wg sync.WaitGroup

	wg.Add(2)
	for j := 0; j < 2; j++ {
		go func(i int) {
			worker(ctx, dataChannel, i)
			wg.Done()
		}(j)
	}
// ...
```

This line of code initializes a waitgroup and runs the goroutines in the background. `sync.WaitGroup` is used to make sure that we are waiting for
goroutines to finish the processing so we can exit the code.

3. 
```go
// ...
	i := 0
	for i < 5 {
		dataChannel <- i
		i++
	}

	wg.Wait()
// ..
```

In this lines of code we are writing simple numbers to the channel, so the goroutines has something to read.

In the last line, we are waiting on the `sync.WaitGroup` which means we are waiting goroutines to finish.
4. 

```go
func worker(ctx context.Context, ch chan int, number int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("context canceled, exiting for worker number %d\n", number)
			return
		case i, ok := <-ch:
			// It means channel is closed
			if !ok {
				log.Printf("channel is closed, exiting worker %d\n", number)
				return
			}
			log.Printf("worker: %d: read %v from channel\n", number, i)
		}
	}
}
```

worker function takes the context parameter and the data channel. 
- ctx is passed to check if the context is canceled or not.
	- if the context is canceled, the goroutine will exit which is the wanted behaviour of the code.
- channel is passed to read the data from the main goroutine
	- if the channel is closed, it means that there are no data to read from so we can exit.


When you run the code, here is what is going to happen:

1. The goroutines will start to listen the channel.
2. The main goroutine will write to the channel.
3. Goroutines will receive some values from the channel and print some stuff.
4. When the 2 second passes, the context will be canceled which means that the goroutines will also exit and the waitgroup value will be 0.
5. the program will exit.

Let's run the code and see the output.

```bash
go run main.go
```

output

```
2023/11/24 21:15:27 worker: 1: read 0 from channel
2023/11/24 21:15:27 worker: 1: read 2 from channel
2023/11/24 21:15:27 worker: 0: read 1 from channel
2023/11/24 21:15:27 worker: 0: read 4 from channel
2023/11/24 21:15:27 worker: 1: read 3 from channel
2023/11/24 21:15:29 context canceled, exiting for worker number 0
2023/11/24 21:15:29 context canceled, exiting for worker number 1
```

As we can see from the output, the goroutines read some values from the channel, then waited for the context to cancel.
When the context is canceled, it is the end time for the goroutines, so they exit, which also results in the main program exit.

From the times of the logs, we can see that the program took nearly 2 seconds to finish.

### Channel Close Case

Let's say before the timeout, the channel is closed by the main goroutine which can mean

- we wrote all of the data we want to process into the channel, it is time to go.
- we processed all of the messages from the channel before the timeout, which is exactly what is wanted in real world case,
the program should finish before the timeout value.

The goroutines should process all of the messages and exit when the channel is closed. 

So, let's modify our code to close the channel when there is no data to write to channel.

{{< highlight go "linenos=table,hl_lines=31" >}}
package main

import (
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dataChannel := make(chan int)

	var wg sync.WaitGroup

	wg.Add(2)
	for j := 0; j < 2; j++ {
		go func(i int) {
			worker(ctx, dataChannel, i)
			wg.Done()
		}(j)
	}

	i := 0
	for i < 5 {
		dataChannel <- i
		i++
	}
	close(dataChannel)

	wg.Wait()
}

func worker(ctx context.Context, ch chan int, number int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("context canceled, exiting for worker number %d\n", number)
			return
		case i, ok := <-ch:
			// It means channel is closed
			if !ok {
				log.Printf("channel is closed, exiting worker %d\n", number)
				return
			}
			log.Printf("worker: %d: read %v from channel\n", number, i)
		}
	}
}
{{< / highlight >}}

At line 31, we are closing the channel.

Now let's run the code again and see the output

```bash
> go run main.go
```

Output

```
2023/11/24 21:21:07 worker: 1: read 0 from channel
2023/11/24 21:21:07 worker: 1: read 2 from channel
2023/11/24 21:21:07 worker: 1: read 3 from channel
2023/11/24 21:21:07 worker: 1: read 4 from channel
2023/11/24 21:21:07 channel is closed, exiting worker 1
2023/11/24 21:21:07 worker: 0: read 1 from channel
2023/11/24 21:21:07 channel is closed, exiting worker 0
```

What we see from the output is that whenever the channel is closed, the worker exits and the program finishes.

You can see from the log times that program exited really quickly before the 2 second timeout of the context, there was no chance for the context to cancel.

Our program is safe by checking
1. context cancellation
2. channel close check

There are 2 conditions we want our program to finish
1. either process all of the data and exit
2. or exit after the timeout.

In golang, the best practice for goroutines is to know when they should exit and passing `context.Context` and checking if it is cancelled is one of 
the best way to handle the goroutines.

In this example we used the `context.Timeout` but there are other options as we defined in the [Context Package Section](#context-package). 

You can use
- [context.WithCancel](https://pkg.go.dev/context#WithCancel) and manually cancel the context whenever you want to.
- [context.WithDeadline](https://pkg.go.dev/context#WithDeadline) and set a deadline for your program to finish.

I hope this blog helped you understand how to use context to propagate the program exit to the sub-goroutines. 

Any feedback is appreciated.

## REFERENCES

- [pkg.go.dev/context](https://pkg.go.dev/context)