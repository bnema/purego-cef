// init.go is the composition root for the cef package. It wires
// the loader, CAPI adapter, and core Engine together.
//
// This file is handwritten — the rest of cef/ is generated.
package cef

import (
	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	"github.com/bnema/purego-cef/internal/loader"
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
	eng = core.New(adapter)
	return eng.Init(settings)
}

// Shutdown releases all CEF resources.
func Shutdown() {
	if eng != nil {
		eng.Shutdown()
	}
}

// DoMessageLoopWork pumps the CEF message loop for one iteration.
func DoMessageLoopWork() {
	if eng != nil {
		eng.DoMessageLoopWork()
	}
}

// MaybeExitSubprocess calls cef_execute_process and exits if this is a subprocess.
func MaybeExitSubprocess() {
	handle, err := loader.Open("")
	if err != nil {
		return
	}
	adapter := capi.New(handle)
	e := core.New(adapter)
	e.MaybeExitSubprocess()
}
