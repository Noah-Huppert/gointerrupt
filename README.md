[![GoDoc](https://godoc.org/github.com/Noah-Huppert/gointerrupt?status.svg)](https://godoc.org/github.com/Noah-Huppert/gointerrupt)
# Go Interrupt
Easily handle interrupts via context.

# Table Of Contents
- [Overview](#overview)
- [Examples](#examples)
  - [HTTP Server](#http-server-example)
  - [Custom Signal](#custom-signal-example)

# Overview
Makes it easy to gracefully shutdown programs.

Provides a graceful shutdown `context.Context` which is canceled when the first 
interrupt signal (`SIGINT`) is received. A harsh shutdown `context.Context` is 
also provided and is canceled when the first termination signal (`SIGTERM`) 
is received.

This allows processes to attempt to gracefully shutdown components by passing
a context. The harsh shutdown signal can be used to kill a graceful shutdown
processes as nicely as possible.

# Examples
## HTTP Server Example
This example shows how to use gointerrupt to gracefully shutdown 
a `net/http.Server`.

```go
package main

import (
	"context"
	"net/http"
	"sync"
	
	"github.com/Noah-Huppert/gointerrupt"
)

func main() {
	// Initialize a go interrupt context pair
	ctxPair := gointerrupt.NewCtxPair(context.Background())

	// Wait group is used to not exit the program until the HTTP server go
	// routine has completed
	var wg sync.WaitGroup

	server := http.Server{
		Addr: ":5000",
	}

	// Run HTTP server
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
		wg.Done()
	}()

	// Gracefully shutdown HTTP server when SIGINT received
	go func() {
		<-ctxPair.Graceful().Done()

		if err := server.Shutdown(ctxPair.Harsh()); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
```

## Custom Signal Example
This example shows how to cancel a context with custom signals.

```go
package main

import (
	"context"
	"net/http"
	"syscall"
	
	"github.com/Noah-Huppert/gointerrupt"
)

func main() {
	// Setup a context to cancel when a kill signal is sent to the process
	ctx := gointerrupt.NewSignalCtx(context.Background(), syscall.SIGKILL)
	
	// Context will cancel when SIGKILL received
	<-ctx.Ctx().Done()
}
```
