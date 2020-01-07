// Easily handles SIGINT and SIGTERM by cancelling a context
package gointerrupt

import (
	"context"
	"os"
	"os/signal"
)

// NewCtx returns a context which will be canceled when an interrupt
// is received.
func NewCtx() (context.Context, context.CancelFunc) {
	ctx, cancelCtx := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		cancelCtx()
	}()

	return ctx, cancelCtx
}
