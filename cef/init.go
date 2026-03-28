// init.go is the composition root for the cef package. It wires
// the loader, CAPI bridge, and core Engine together.
//
// This file is handwritten — the rest of cef/ is generated.
//
// The hexagonal refactor removed the following hand-written APIs that are
// now provided (or intentionally dropped) by the generated public layer:
//   - NewKeyEvent, KeyEventSetType → use generated KeyEvent struct directly
//   - DefaultWindowInfo, DefaultBrowserSettings → use generated zero-value structs
//   - SetHandlerTraceEnabled, HandlerTraceEnabled → handler tracing removed
package cef

import (
	"fmt"
	"os"
	"runtime"
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
// It runs once via sync.Once — if initialization fails, subsequent calls return
// the cached error without retrying. This is by design: a failed CEF init leaves
// the process in an indeterminate state, so retrying would be unsafe.
func Init(settings Settings) error {
	return InitWithApp(settings, nil)
}

// InitWithApp loads the CEF library, registers all symbols, and initializes the
// runtime with a custom App handler. The App enables registration of custom
// schemes, browser process handlers, and subprocess apps.
// It runs once via sync.Once — if initialization fails, subsequent calls return
// the cached error without retrying.
func InitWithApp(settings Settings, app App) error {
	initOnce.Do(func() {
		handle, err := loader.Open(settings.CEFDir)
		if err != nil {
			initErr = err
			return
		}
		bridge := capi.NewBridge(handle)
		e := core.New(bridge)
		eng = e

		var appPtr unsafe.Pointer
		var wrapped App
		if app != nil {
			wrapped = NewApp(app)
			appPtr = extractRawPointer(wrapped)
		}
		initErr = e.InitWithApp(settings, appPtr)
		runtime.KeepAlive(wrapped)
	})
	return initErr
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

// MaybeExitSubprocess uses a lightweight path that binds only
// cef_execute_process, bypassing the full Bridge initialization.
// This is intentional — subprocess detection must happen before
// the full CEF runtime is initialized.
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
	runtime.KeepAlive(args)
	if code >= 0 {
		os.Exit(int(code))
	}
}
