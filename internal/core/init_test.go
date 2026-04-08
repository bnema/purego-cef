package core

import (
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/stretchr/testify/mock"

	"github.com/bnema/purego-cef/internal/ports/out/mocks"
)

// readCEFString extracts a Go string from a CEFStringT for test assertions.
func readCEFString(cs CEFStringT) string {
	if cs.Str == nil || cs.Length == 0 {
		return ""
	}
	return string(utf16.Decode(unsafe.Slice(cs.Str, cs.Length)))
}

// newTestEngine creates an Engine backed by a MockCAPI whose StringSet
// copies UTF-16 data into the output CEFStringT (mirroring real CEF
// behavior) and whose StringClear / NewCallback are safe no-ops.
func newTestEngine(t *testing.T) *Engine {
	m := mocks.NewMockCAPI(t)

	// StringSet: copy UTF-16 data into the output CEFStringT so field
	// assertions work without a real CEF library.
	m.EXPECT().StringSet(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(src *uint16, srcLen uintptr, output unsafe.Pointer, _ int32) int32 {
			out := (*CEFStringT)(output)
			n := int(srcLen)
			if n == 0 || src == nil {
				out.Str = nil
				out.Length = 0
				return 1
			}
			buf := make([]uint16, n)
			copy(buf, unsafe.Slice(src, n))
			out.Str = &buf[0]
			out.Length = uintptr(n)
			return 1
		}).Maybe()

	// StringClear: no-op (no real allocator to free).
	m.EXPECT().StringClear(mock.Anything).Maybe()

	// NewCallback: return zero (unused by settingsToRaw).
	m.EXPECT().NewCallback(mock.Anything).Maybe().Return(uintptr(0))

	return New(m)
}

func TestSettingsToRaw_CachePaths(t *testing.T) {
	e := newTestEngine(t)

	s := Settings{
		CachePath:     "/tmp/cef-cache",
		RootCachePath: "/tmp/cef-root-cache",
	}

	raw, cleanup := e.settingsToRaw(s)
	defer cleanup()

	if got := readCEFString(raw.CachePath); got != "/tmp/cef-cache" {
		t.Errorf("CachePath = %q, want %q", got, "/tmp/cef-cache")
	}
	if got := readCEFString(raw.RootCachePath); got != "/tmp/cef-root-cache" {
		t.Errorf("RootCachePath = %q, want %q", got, "/tmp/cef-root-cache")
	}
}

func TestSettingsToRaw_CachePathsEmpty(t *testing.T) {
	e := newTestEngine(t)

	s := Settings{} // no cache paths set

	raw, cleanup := e.settingsToRaw(s)
	defer cleanup()

	if raw.CachePath.Str != nil {
		t.Errorf("CachePath.Str should be nil when not set")
	}
	if raw.RootCachePath.Str != nil {
		t.Errorf("RootCachePath.Str should be nil when not set")
	}
}

func TestSettingsToRaw_AllStringFields(t *testing.T) {
	e := newTestEngine(t)

	s := Settings{
		BrowserSubprocessPath: "/usr/lib/cef/subprocess",
		LogFile:               "/var/log/cef.log",
		CachePath:             "/cache",
		RootCachePath:         "/root-cache",
	}

	raw, cleanup := e.settingsToRaw(s)
	defer cleanup()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"BrowserSubprocessPath", readCEFString(raw.BrowserSubprocessPath), "/usr/lib/cef/subprocess"},
		{"LogFile", readCEFString(raw.LogFile), "/var/log/cef.log"},
		{"CachePath", readCEFString(raw.CachePath), "/cache"},
		{"RootCachePath", readCEFString(raw.RootCachePath), "/root-cache"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}
