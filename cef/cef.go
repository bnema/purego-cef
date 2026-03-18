// Package cef provides the public API for embedding the Chromium Embedded
// Framework via purego (no cgo).
package cef

import (
	"fmt"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/cefstr"
	"github.com/bnema/purego-cef/internal/loader"
)

// Init opens the CEF shared library from runtimeDir, registers all bindings,
// and calls cef_initialize. An empty runtimeDir uses the default location
// (~/.local/share/cef or $CEF_DIR).
func Init(runtimeDir string, settings Settings) error {
	handle, err := loader.Open(runtimeDir)
	if err != nil {
		return fmt.Errorf("cef: open loader: %w", err)
	}

	capi.Register(handle)
	cefstr.Bind(handle)

	args := NewMainArgsFromOS()

	cs, cleanup, err := settings.toC()
	if err != nil {
		return fmt.Errorf("cef: convert settings: %w", err)
	}
	defer cleanup()

	ok := capi.CEFInitialize(args.Ptr(), unsafe.Pointer(&cs), nil, nil)
	if ok != 1 {
		return fmt.Errorf("cef: cef_initialize returned %d", ok)
	}
	return nil
}

// Shutdown releases all CEF resources. Must be called on the main thread
// before the process exits.
func Shutdown() {
	capi.CEFShutdown()
}

// DoMessageLoopWork performs a single iteration of CEF message loop work.
// Use this when running with external_message_pump enabled.
func DoMessageLoopWork() {
	capi.CEFDoMessageLoopWork()
}

// MaybeExitSubprocess calls cef_execute_process and returns the exit code.
// If the exit code is >= 0, the caller should os.Exit with that code (this
// process was a CEF sub-process). If < 0, continue as the browser process.
func MaybeExitSubprocess(runtimeDir string) (int, error) {
	handle, err := loader.Open(runtimeDir)
	if err != nil {
		return -1, fmt.Errorf("cef: open loader: %w", err)
	}

	capi.Register(handle)
	cefstr.Bind(handle)

	args := NewMainArgsFromOS()
	code := capi.CEFExecuteProcess(args.Ptr(), nil, nil)
	return int(code), nil
}
