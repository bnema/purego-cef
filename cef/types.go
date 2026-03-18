package cef

import (
	"os"
	"unsafe"

	"github.com/bnema/purego-cef/internal/capi"
	"github.com/bnema/purego-cef/internal/cefstr"
)

// MainArgs holds the command-line arguments passed to CEF.
// The backing byte buffers are kept alive for the lifetime of the value.
type MainArgs struct {
	raw  capi.CEFMainArgsT
	argv [][]byte // keep backing buffers alive
	ptrs []*byte  // keep pointer array alive
}

// NewMainArgs creates MainArgs from a string slice.
func NewMainArgs(args []string) MainArgs {
	var m MainArgs
	m.argv = make([][]byte, len(args))
	m.ptrs = make([]*byte, len(args))
	for i, s := range args {
		// Null-terminated C string.
		buf := make([]byte, len(s)+1)
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

// NewMainArgsFromOS creates MainArgs from os.Args.
func NewMainArgsFromOS() MainArgs { return NewMainArgs(os.Args) }

// Ptr returns an unsafe.Pointer to the underlying C struct.
func (m *MainArgs) Ptr() unsafe.Pointer { return unsafe.Pointer(&m.raw) }

// Settings configures the CEF runtime.
type Settings struct {
	BrowserSubprocessPath      string
	ExternalMessagePump        bool
	WindowlessRenderingEnabled bool
	NoSandbox                  bool
	LogSeverity                int32
	DisableSignalHandlers      bool
}

// DefaultSettings returns settings suitable for off-screen rendering with
// an external message pump.
func DefaultSettings() Settings {
	return Settings{
		ExternalMessagePump:        true,
		WindowlessRenderingEnabled: true,
		NoSandbox:                  true,
		DisableSignalHandlers:      true,
	}
}

// toC converts the high-level Settings to the C struct.
// The returned cleanup function must be called after CEFInitialize returns.
func (s Settings) toC() (capi.CEFSettingsT, func(), error) {
	var c capi.CEFSettingsT
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
	if s.DisableSignalHandlers {
		c.DisableSignalHandlers = 1
	}
	c.LogSeverity = capi.CEFLogSeverityT(s.LogSeverity)

	var cleanups []func()
	cleanup := func() {
		for _, fn := range cleanups {
			fn()
		}
	}

	if s.BrowserSubprocessPath != "" {
		str, cl, err := cefstr.FromGo(s.BrowserSubprocessPath)
		if err != nil {
			return c, nil, err
		}
		cleanups = append(cleanups, cl)
		c.BrowserSubprocessPath = capi.CEFStringT{
			Str:    str.Str,
			Length: str.Length,
			Dtor:   str.Dtor,
		}
	}

	return c, cleanup, nil
}

// Rect represents a rectangle.
type Rect struct {
	X, Y, W, H int32
}
