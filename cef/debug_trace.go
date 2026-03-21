package cef

import (
	"fmt"
	"os"
	"sync/atomic"
)

var (
	handlerTraceEnabled atomic.Bool
	handlerTraceSeq     atomic.Uint64
)

// SetHandlerTraceEnabled toggles diagnostic tracing for handler wrappers and
// tracked refcount events.
func SetHandlerTraceEnabled(enabled bool) {
	handlerTraceEnabled.Store(enabled)
	if enabled {
		handlerTraceSeq.Store(0)
	}
}

// HandlerTraceEnabled reports whether handler tracing is currently enabled.
func HandlerTraceEnabled() bool {
	return handlerTraceEnabled.Load()
}

func traceHandlerf(format string, args ...any) {
	if !handlerTraceEnabled.Load() {
		return
	}
	seq := handlerTraceSeq.Add(1)
	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(os.Stderr, "[purego-cef:%06d] %s\n", seq, msg)
}
