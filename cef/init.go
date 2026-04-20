// init.go is the composition root for the cef package. It wires
// the loader, CAPI bridge, and core Engine together.
//
// This file is handwritten — the rest of cef/ is generated.
//
// The hexagonal refactor removed the following hand-written APIs that are
// now provided (or intentionally dropped) by the generated public layer:
//   - NewKeyEvent, KeyEventSetType → use generated KeyEvent struct directly
//   - DefaultWindowInfo, DefaultBrowserSettings → use NewWindowInfo(), NewBrowserSettings()
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
)

// Settings configures the CEF runtime.
//
// This is the user-facing settings type for Init/InitWithApp. Prefer it over
// the raw generated CEFSettings struct unless you specifically need the exact
// CEF memory layout.
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
// cef_execute_process, bypassing the full runtime initialization.
// This is intentional — subprocess detection must happen before
// the full CEF runtime is initialized.
func MaybeExitSubprocess() {
	MaybeExitSubprocessWithApp(nil)
}

// MaybeExitSubprocessWithApp is like MaybeExitSubprocess but passes the given
// App to cef_execute_process. This enables custom render/browser-process
// handlers in helper subprocesses without calling cef_initialize first.
func MaybeExitSubprocessWithApp(app App) {
	handle, err := loader.Open("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cef: MaybeExitSubprocessWithApp: %v\n", err)
		return
	}

	// Wrapping App handlers uses the generated refcount helpers, which depend on
	// the package-level engine reference manager. Seed it with a lightweight core
	// engine backed by the raw CAPI bridge just for the duration of execute_process.
	prevEng := eng
	bridge := capi.NewBridge(handle)
	eng = core.New(bridge)
	defer func() {
		eng = prevEng
	}()

	var appPtr unsafe.Pointer
	var wrapped App
	if app != nil {
		wrapped = NewApp(app)
		appPtr = extractRawPointer(wrapped)
	}

	args := core.NewMainArgs(os.Args)
	code := bridge.ExecuteProcess(args.Ptr(), appPtr, nil)
	runtime.KeepAlive(args)
	runtime.KeepAlive(wrapped)
	if code >= 0 {
		os.Exit(int(code))
	}
}
