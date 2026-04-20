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
	"io"
	"os"
	"runtime"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/core"
	"github.com/bnema/purego-cef/internal/loader"
	portout "github.com/bnema/purego-cef/internal/ports/out"
)

var (
	openCEFLibrary           = loader.Open
	newCAPI                  = func(handle uintptr) portout.CAPI { return capi.NewBridge(handle) }
	processArgs              = func() []string { return os.Args }
	exitProcess              = os.Exit
	stderrWriter   io.Writer = os.Stderr
)

// Settings configures the CEF runtime.
//
// This is the user-facing settings type for Init/InitWithApp. Prefer it over
// the raw generated RawSettings struct unless you specifically need the exact
// CEF memory layout.
type Settings struct {
	CEFDir                     string
	LogSeverity                int32
	MultiThreadedMessageLoop   bool
	WindowlessRenderingEnabled bool
	ExternalMessagePump        bool
	NoSandbox                  bool
	BrowserSubprocessPath      string
	LogFile                    string
	// Deprecated: reserved for bootstrap diagnostics only; it is not translated
	// into a CEF setting and has no effect in purego-cef itself.
	InitTraceFile string
	CachePath     string
	RootCachePath string
}

func (s Settings) coreSettings() core.Settings {
	return core.Settings{
		CEFDir:                     s.CEFDir,
		LogSeverity:                s.LogSeverity,
		MultiThreadedMessageLoop:   s.MultiThreadedMessageLoop,
		WindowlessRenderingEnabled: s.WindowlessRenderingEnabled,
		ExternalMessagePump:        s.ExternalMessagePump,
		NoSandbox:                  s.NoSandbox,
		BrowserSubprocessPath:      s.BrowserSubprocessPath,
		LogFile:                    s.LogFile,
		InitTraceFile:              s.InitTraceFile,
		CachePath:                  s.CachePath,
		RootCachePath:              s.RootCachePath,
	}
}

// DefaultSettings returns Settings suitable for off-screen rendering.
func DefaultSettings() Settings {
	return Settings{
		ExternalMessagePump:        true,
		WindowlessRenderingEnabled: true,
		NoSandbox:                  true,
	}
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
		handle, err := openCEFLibrary(settings.CEFDir)
		if err != nil {
			initErr = err
			return
		}
		capi := newCAPI(handle)
		e := core.New(capi)
		eng = e

		var appPtr unsafe.Pointer
		var wrapped App
		if app != nil {
			wrapped = NewApp(app)
			appPtr = extractRawPointer(wrapped)
		}
		initErr = e.InitWithApp(settings.coreSettings(), appPtr)
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

// ExecuteSubprocess runs cef_execute_process using a lightweight path that
// binds only the symbols needed for subprocess detection.
//
// It returns executed=true when the current process was handled as a CEF helper
// subprocess. In that case exitCode should be used as the process exit status.
// When executed=false and err=nil, the caller should continue normal startup.
func ExecuteSubprocess() (executed bool, exitCode int, err error) {
	handle, err := openCEFLibrary("")
	if err != nil {
		return false, 0, err
	}
	capi := newCAPI(handle)
	args := core.NewMainArgs(processArgs())
	code := capi.ExecuteProcess(args.Ptr(), nil, nil)
	runtime.KeepAlive(args)
	if code >= 0 {
		return true, int(code), nil
	}
	return false, 0, nil
}

// MaybeExitSubprocess is a convenience helper for main packages.
// Library code should usually prefer ExecuteSubprocess.
func MaybeExitSubprocess() {
	executed, exitCode, err := ExecuteSubprocess()
	if err != nil {
		fmt.Fprintf(stderrWriter, "cef: ExecuteSubprocess: %v\n", err)
		return
	}
	if executed {
		exitProcess(exitCode)
	}
}

// MaybeExitSubprocessWithApp is like MaybeExitSubprocess but passes the given
// App to cef_execute_process. This enables custom render/browser-process
// handlers in helper subprocesses without calling cef_initialize first.
//
// This is a convenience helper for main packages. It may write to stderr and
// call os.Exit. Prefer ExecuteSubprocess when you do not need a custom App.
func MaybeExitSubprocessWithApp(app App) {
	if app == nil {
		MaybeExitSubprocess()
		return
	}

	handle, err := openCEFLibrary("")
	if err != nil {
		fmt.Fprintf(stderrWriter, "cef: MaybeExitSubprocessWithApp: %v\n", err)
		return
	}

	// Wrapping App handlers uses the generated refcount helpers, which depend on
	// the package-level engine reference manager. Seed it with a lightweight core
	// engine backed by the raw CAPI bridge just for the duration of execute_process.
	prevEng := eng
	capi := newCAPI(handle)
	eng = core.New(capi)
	defer func() {
		eng = prevEng
	}()

	var appPtr unsafe.Pointer
	var wrapped App
	if app != nil {
		wrapped = NewApp(app)
		appPtr = extractRawPointer(wrapped)
	}

	args := core.NewMainArgs(processArgs())
	code := capi.ExecuteProcess(args.Ptr(), appPtr, nil)
	runtime.KeepAlive(args)
	runtime.KeepAlive(wrapped)
	if code >= 0 {
		exitProcess(int(code))
	}
}
