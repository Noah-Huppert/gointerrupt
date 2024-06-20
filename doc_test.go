package gointerrupt

import (
	"context"
	"log"
	"net/http"
	"sync"
	"syscall"
	"time"
)

// Demonstrates how a CtxPair can be used to gracefully shutdown
// a net/http.Server.
func ExampleCtxPair() {
	// Initialize a go interrupt context pair
	ctxPair := NewCtxPair(context.Background())

	// Wait group is used to not exit the program until the HTTP server go
	// routine has completed
	var wg sync.WaitGroup

	server := http.Server{
		Addr: ":5000",
	}

	// Run HTTP server
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
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

// Shows how to setup a context.Context to cancel when custom signals are
// received by the process.
func ExampleSignalCtx() {
	// Setup a context to cancel when a kill signal is sent to the process
	ctx := NewSignalCtx(context.Background(), syscall.SIGKILL)

	// Context will cancel when SIGKILL received
	<-ctx.Ctx().Done()
}

// Shows how to manually cancel a context.
func ExampleSignalCtx_Cancel() {
	// Setup a context which will cancel when the SIGKILL signal is received
	ctx := NewSignalCtx(context.Background(), syscall.SIGKILL)

	// Create a timer which will trigger in 10 seconds
	timer := time.NewTimer(time.Second * 10)

	go func() {
		<-timer.C
		log.Println("Timer went off! Cancelling context")
		ctx.Cancel()
	}()

	<-ctx.Ctx().Done()
	log.Println("Context was canceled")
}
