// init.go is the composition root for the cef package. It wires
// the loader, CAPI adapter, and core Engine together.
//
// This file is handwritten — the rest of cef/ is generated.
package cef

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	"github.com/bnema/purego-cef/internal/loader"
	"github.com/ebitengine/purego"
)

// Settings configures the CEF runtime.
type Settings = core.Settings

// DefaultSettings returns Settings suitable for off-screen rendering.
func DefaultSettings() Settings {
	return core.DefaultSettings()
}

// Init loads the CEF library, registers all symbols, and initializes the runtime.
func Init(settings Settings) error {
	handle, err := loader.Open(settings.CEFDir)
	if err != nil {
		return err
	}
	adapter := capi.New(handle)
	e := core.New(adapter)
	engPtr.Store(e)
	return e.Init(settings)
}

// Shutdown releases all CEF resources.
func Shutdown() {
	if e := engPtr.Load(); e != nil {
		e.Shutdown()
	}
}

// DoMessageLoopWork pumps the CEF message loop for one iteration.
func DoMessageLoopWork() {
	if e := engPtr.Load(); e != nil {
		e.DoMessageLoopWork()
	}
}

// MaybeExitSubprocess calls cef_execute_process and exits if this is a subprocess.
// Uses a lightweight path that only binds cef_execute_process.
func MaybeExitSubprocess() {
	handle, err := loader.Open("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cef: MaybeExitSubprocess: %v\n", err)
		return
	}
	var executeProcess func(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer) int32
	purego.RegisterLibFunc(&executeProcess, handle, "cef_execute_process")
	args := core.NewMainArgs(os.Args)
	code := executeProcess(args.Ptr(), nil, nil)
	if code >= 0 {
		os.Exit(int(code))
	}
}
