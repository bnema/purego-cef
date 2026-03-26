// internal/core/engine.go
package core

import "unsafe"

// CAPI is the outbound port interface. Core depends on this for all C
// interactions. The real implementation (internal/capi) uses purego;
// tests inject mocks.
//
// This interface will grow as we migrate functionality from cef/.
// Eventually it will be replaced by generated interfaces in ports/out.
type CAPI interface {
	// Functions — CEF free functions
	Initialize(args, settings, app, sandboxInfo unsafe.Pointer) int32
	Shutdown()
	DoMessageLoopWork()
	ExecuteProcess(args, app, sandboxInfo unsafe.Pointer) int32

	// Callbacks — wraps purego.NewCallback
	NewCallback(fn any) uintptr

	// Strings — CEF UTF-16 string operations
	StringSet(src *uint16, srcLen uintptr, output unsafe.Pointer, copy int32) int32
	StringClear(s unsafe.Pointer)
	StringUserfreeFree(s unsafe.Pointer)
	StringListSize(list uintptr) uintptr
	StringListValue(list uintptr, index uintptr, value unsafe.Pointer) int32

	// Registration — bind all symbols from handle
	Register(handle uintptr)
}

// Engine is the core domain. It holds no global state — all state
// flows through the instance. The cef/ composition root creates one
// Engine at Init() time and holds it as a package-level variable.
type Engine struct {
	capi        CAPI
	refs        *RefManager
	initialized bool
}

// New creates an Engine with the given CAPI adapter.
func New(capi CAPI) *Engine {
	e := &Engine{capi: capi}
	e.refs = NewRefManager(capi)
	return e
}

// Refs returns the engine's RefManager.
func (e *Engine) Refs() *RefManager {
	return e.refs
}
