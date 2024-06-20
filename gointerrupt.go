// Package gointerrupt makes it easily to handle signals by canceling a context.
//
// Use CtxPair to handle graceful shutdown of a program. See the example for a
// familiar net/http.Server setup.
//
// Use SignalCtx to cancel contexts when custom signals are received.
package gointerrupt

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// CtxPair provides 2 contexts which can be used to exit a program gracefully
// when shutdown signals are received by the process.
//
// The graceful context will canceled when an interrupt signal is received.
// The harsh context will be canceled when a terminate signal is received.
//
// The graceful context being canceled indicates graceful shutdown should begin,
// and the harsh indicates when it should end.
type CtxPair struct {
	// gracefullCtx will be canceled when a signal is received which indicates
	// graceful shutdown should begin.
	gracefullCtx SignalCtx

	// harshCtx will be canceled when a signal is received which indicates
	// graceful shutdown should end.
	harshCtx SignalCtx
}

// NewCtxPair creates a new CtxPair which cancels the graceful context on SIGINT
// and harsh context on SIGTERM.
func NewCtxPair(bkgrnd context.Context) CtxPair {
	return CtxPair{
		gracefullCtx: NewSignalCtx(bkgrnd, syscall.SIGINT),
		harshCtx:     NewSignalCtx(bkgrnd, syscall.SIGTERM),
	}
}

// Graceful returns the context which is canceled when graceful shutdown
// should begin.
func (pair CtxPair) Graceful() context.Context {
	return pair.gracefullCtx.Ctx()
}

// Harsh returns the context which is canceled when graceful shutdown
// should end.
func (pair CtxPair) Harsh() context.Context {
	return pair.harshCtx.Ctx()
}

// SignalCtx is a context which cancels when a signal is received by the process
type SignalCtx struct {
	// ctx is the context to be canceled
	ctx context.Context

	// cancelFn is the function which when called will cancel the context
	cancelFn func()
}

// NewSignalCtx creates a SignalCtx for the specified signal
func NewSignalCtx(bkgrnd context.Context, sig os.Signal) SignalCtx {
	ctx, cancelCtx := context.WithCancel(bkgrnd)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sig)

	go func() {
		<-sigs
		cancelCtx()
	}()

	return SignalCtx{
		ctx:      ctx,
		cancelFn: cancelCtx,
	}
}

// Ctx returns the context which is canceled when a signal is received.
func (sigCtx SignalCtx) Ctx() context.Context {
	return sigCtx.ctx
}

// Cancel forcefully cancels the context.
func (sigCtx SignalCtx) Cancel() {
	sigCtx.cancelFn()
}
