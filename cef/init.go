package cef

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"unsafe"

	"github.com/bnema/purego-cef/cef/internal/raw"
	"github.com/ebitengine/purego"
)

// ---------------------------------------------------------------------------
// Settings
// ---------------------------------------------------------------------------

// Settings configures the CEF runtime.
type Settings struct {
	// CEFDir is the path to the CEF runtime directory containing libcef.so.
	// An empty string resolves via $CEF_DIR or ~/.local/share/cef.
	CEFDir string

	// LogSeverity controls the CEF log level.
	LogSeverity int32

	// WindowlessRenderingEnabled enables off-screen rendering mode.
	WindowlessRenderingEnabled bool

	// ExternalMessagePump enables the external message pump model.
	ExternalMessagePump bool

	// NoSandbox disables the CEF sandbox.
	NoSandbox bool
}

// DefaultSettings returns Settings suitable for off-screen rendering with
// an external message pump.
func DefaultSettings() Settings {
	return Settings{
		ExternalMessagePump:        true,
		WindowlessRenderingEnabled: true,
		NoSandbox:                  true,
	}
}

// DefaultWindowInfo returns a WindowInfo with the Size field correctly
// initialized. Callers should set WindowlessRenderingEnabled and other
// fields as needed after construction.
func DefaultWindowInfo() WindowInfo {
	var info WindowInfo
	info.Size = unsafe.Sizeof(info)
	return info
}

// DefaultBrowserSettings returns a BrowserSettings with the Size field
// correctly initialized. Callers may override individual fields as needed.
func DefaultBrowserSettings() BrowserSettings {
	var s BrowserSettings
	s.Size = unsafe.Sizeof(s)
	return s
}

// toRaw converts Settings to the raw C struct. The returned cleanup function
// must be called after cef_initialize returns.
func (s Settings) toRaw() (raw.CEFSettingsT, func()) {
	var c raw.CEFSettingsT
	c.Size = uintptr(unsafe.Sizeof(c))

	if s.NoSandbox {
		c.NoSandbox = 1
	}
	if s.ExternalMessagePump {
		c.ExternalMessagePump = 1
	}
	if s.WindowlessRenderingEnabled {
		c.WindowlessRenderingEnabled = 1
	}
	c.LogSeverity = raw.CEFLogSeverityT(s.LogSeverity)

	return c, func() {} // cleanup placeholder; strings would be freed here
}

// ---------------------------------------------------------------------------
// MainArgs
// ---------------------------------------------------------------------------

// MainArgs holds the command-line arguments passed to CEF. The backing byte
// buffers are kept alive for the lifetime of the value.
type MainArgs struct {
	raw  raw.CEFMainArgsT
	argv [][]byte // keep backing buffers alive
	ptrs []*byte  // keep pointer array alive
}

// NewMainArgs creates MainArgs from a string slice.
func NewMainArgs(args []string) *MainArgs {
	m := &MainArgs{}
	m.argv = make([][]byte, len(args))
	m.ptrs = make([]*byte, len(args))
	for i, s := range args {
		buf := make([]byte, len(s)+1) // null-terminated C string
		copy(buf, s)
		m.argv[i] = buf
		m.ptrs[i] = &buf[0]
	}
	m.raw.Argc = int32(len(args))
	if len(m.ptrs) > 0 {
		m.raw.Argv = uintptr(unsafe.Pointer(&m.ptrs[0]))
	}
	return m
}

// Ptr returns an unsafe.Pointer to the underlying C struct.
func (m *MainArgs) Ptr() unsafe.Pointer { return unsafe.Pointer(&m.raw) }

// ---------------------------------------------------------------------------
// Library loading & version validation (ported from internal/loader)
// ---------------------------------------------------------------------------

const defaultCEFVersion = 145

// cefVersionInfoChromeMajor is the entry index for the chrome major version
// in cef_version_info().
const cefVersionInfoChromeMajor = 4

func resolveDir(arg string) (string, error) {
	if env := os.Getenv("CEF_DIR"); env != "" {
		return env, nil
	}
	if arg != "" {
		return arg, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".local", "share", "cef"), nil
}

func targetMajor() int32 {
	if raw := os.Getenv("CEF_VERSION"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			return int32(n)
		}
	}
	return defaultCEFVersion
}

func openLib(dir string) (uintptr, error) {
	runtimeDir, err := resolveDir(dir)
	if err != nil {
		return 0, err
	}
	libPath := filepath.Join(runtimeDir, "libcef.so")
	handle, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, fmt.Errorf("dlopen %s: %w", libPath, err)
	}

	// Validate version unless explicitly skipped.
	if os.Getenv("CEF_SKIP_VERSION_CHECK") != "1" {
		sym, err := purego.Dlsym(handle, "cef_version_info")
		if err != nil {
			return 0, fmt.Errorf("resolve cef_version_info: %w", err)
		}
		var versionInfo func(int32) int32
		purego.RegisterFunc(&versionInfo, sym)
		got := versionInfo(cefVersionInfoChromeMajor)
		want := targetMajor()
		if got != want {
			return 0, fmt.Errorf("unsupported CEF runtime: chrome major=%d want=%d", got, want)
		}
	}
	return handle, nil
}

// ---------------------------------------------------------------------------
// Init / Shutdown / orchestration
// ---------------------------------------------------------------------------

var (
	libHandle   uintptr
	initialized bool
	libOnce     sync.Once
	libErr      error
	initOnce    sync.Once
	initErr     error
)

// loadLibrary opens the CEF shared library and registers all raw symbols.
// It is idempotent — safe to call from both MaybeExitSubprocess and Init.
func loadLibrary(cefDir string) error {
	libOnce.Do(func() {
		handle, err := openLib(cefDir)
		if err != nil {
			libErr = fmt.Errorf("cef: open library: %w", err)
			return
		}
		libHandle = handle
		bindStringFuncs(handle)
		raw.Register(handle)
	})
	return libErr
}

// Init loads the CEF shared library, binds string and raw symbol functions,
// and calls cef_initialize. It is safe to call multiple times; only the first
// call takes effect.
func Init(settings Settings) error {
	initOnce.Do(func() {
		initErr = doInit(settings)
	})
	return initErr
}

func doInit(settings Settings) error {
	if err := loadLibrary(settings.CEFDir); err != nil {
		return err
	}

	// Build main args from os.Args.
	args := NewMainArgs(os.Args)

	// Convert settings to raw struct.
	cs, cleanup := settings.toRaw()
	defer cleanup()

	ok := raw.CEFInitialize(args.Ptr(), unsafe.Pointer(&cs), nil, nil)
	if ok != 1 {
		return fmt.Errorf("cef: cef_initialize returned %d", ok)
	}
	initialized = true
	return nil
}

// Shutdown releases all CEF resources. The library handle is intentionally
// not closed via Dlclose because CEF does not support clean unloading.
func Shutdown() {
	if !initialized {
		return
	}
	raw.CEFShutdown()
	initialized = false
}

// DoMessageLoopWork pumps the CEF message loop for one iteration. Call this
// from your application's main loop when using an external message pump.
func DoMessageLoopWork() {
	raw.CEFDoMessageLoopWork()
}

// MaybeExitSubprocess calls cef_execute_process and exits the current process
// if it was launched as a CEF sub-process. For the browser process this is a
// no-op (returns without exiting).
func MaybeExitSubprocess() {
	if err := loadLibrary(""); err != nil {
		// If the library cannot be loaded there is nothing useful we can do
		// in a sub-process, so bail out with a diagnostic.
		fmt.Fprintf(os.Stderr, "cef: MaybeExitSubprocess: %v\n", err)
		os.Exit(1)
	}
	args := NewMainArgs(os.Args)
	code := raw.CEFExecuteProcess(args.Ptr(), nil, nil)
	if code >= 0 {
		os.Exit(int(code))
	}
}
