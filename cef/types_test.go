package cef

import (
	"testing"
	"unsafe"
)

func TestNewMainArgsFromSlice(t *testing.T) {
	args := NewMainArgs([]string{"app", "--flag"})
	if args.raw.Argc != 2 {
		t.Fatalf("argc = %d, want 2", args.raw.Argc)
	}
	if args.raw.Argv == 0 {
		t.Fatal("argv is zero")
	}
	if unsafe.Sizeof(args.raw) == 0 {
		t.Fatal("unexpected zero-sized args")
	}
}

func TestNewMainArgsEmpty(t *testing.T) {
	args := NewMainArgs(nil)
	if args.raw.Argc != 0 {
		t.Fatalf("argc = %d, want 0", args.raw.Argc)
	}
}

func TestDefaultSettingsEnableOSR(t *testing.T) {
	settings := DefaultSettings()
	if !settings.WindowlessRenderingEnabled {
		t.Fatal("windowless rendering should be enabled")
	}
	if !settings.ExternalMessagePump {
		t.Fatal("external message pump should be enabled")
	}
	if !settings.NoSandbox {
		t.Fatal("no_sandbox should be enabled")
	}
}
