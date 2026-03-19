package cef

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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

// defaultAPIVersion is the CEF API version we target. This must match
// CEF_API_VERSION_LAST from the headers used for generation.
const defaultAPIVersion = 14500

// setAPIVersion patches the internal g_version global inside libcef.so so
// that cef_api_version() returns our target version instead of -1. Without
// this, CEF's versioned CToCpp wrappers reject all client structs.
//
// The cef_api_version function is a trivial global read:
//
//	push %rbp; mov %rsp,%rbp; mov <offset>(%rip),%eax; pop %rbp; ret
//
// We decode the RIP-relative offset to locate g_version and write to it.
func setAPIVersion(handle uintptr) error {
	sym, err := purego.Dlsym(handle, "cef_api_version")
	if err != nil {
		return fmt.Errorf("resolve cef_api_version: %w", err)
	}

	// Decode: mov offset(%rip),%eax at sym+4 (opcode 8b 05 xx xx xx xx)
	// The offset is a 4-byte signed integer at sym+6.
	// RIP points to the next instruction at sym+10.
	offset := *(*int32)(unsafe.Pointer(sym + 6))
	gVersionAddr := sym + 10 + uintptr(offset)

	*(*int32)(unsafe.Pointer(gVersionAddr)) = defaultAPIVersion
	return nil
}

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
		if err := setAPIVersion(handle); err != nil {
			libErr = fmt.Errorf("cef: %w", err)
			return
		}
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

// isSubprocess returns true if os.Args contains a --type= flag, indicating
// this process was spawned by CEF as a renderer, GPU, or utility subprocess.
func isSubprocess() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--type=") {
			return true
		}
	}
	return false
}

// closeGoRuntimeFDs closes all file descriptors above stderr (fd > 2) that
// the Go runtime opened before main(). CEF subprocess FD tracking expects a
// clean process; Go's epoll fd, netpoller pipes, and signal fds trigger a
// fatal "FD ownership violation" in the GPU and renderer processes.
func closeGoRuntimeFDs() {
	entries, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return
	}
	for _, e := range entries {
		fd, err := strconv.Atoi(e.Name())
		if err != nil || fd <= 2 {
			continue
		}
		syscall.Close(fd)
	}
}

// MaybeExitSubprocess calls cef_execute_process and exits the current process
// if it was launched as a CEF sub-process. For the browser process this is a
// no-op (returns without exiting).
//
// When running as a subprocess (--type= in args), all file descriptors above
// stderr are closed before calling CEF to avoid FD ownership violations from
// Go runtime descriptors (epoll, signal pipes, netpoller).
func MaybeExitSubprocess() {
	if err := loadLibrary(""); err != nil {
		// If the library cannot be loaded there is nothing useful we can do
		// in a sub-process, so bail out with a diagnostic.
		fmt.Fprintf(os.Stderr, "cef: MaybeExitSubprocess: %v\n", err)
		os.Exit(1)
	}
	if isSubprocess() {
		closeGoRuntimeFDs()
	}
	args := NewMainArgs(os.Args)
	code := raw.CEFExecuteProcess(args.Ptr(), nil, nil)
	if code >= 0 {
		os.Exit(int(code))
	}
}
