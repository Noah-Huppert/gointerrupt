package gointerrupt

// ExampleNewCtx shows how a new context can be created which will be canceled
// when an interrupt is received.
func ExampleNewCtx() {
	ctx, _ := NewCtx()
	<-ctx.Done()
}
