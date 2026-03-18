package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDirPrecedence(t *testing.T) {
	t.Setenv("CEF_DIR", "/env/cef")
	got, err := resolveDir("/arg/cef", func() (string, error) { return "/home/brice", nil })
	if err != nil {
		t.Fatal(err)
	}
	if got != "/env/cef" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveDirFallsBackToArg(t *testing.T) {
	t.Setenv("CEF_DIR", "")
	got, err := resolveDir("/arg/cef", func() (string, error) { return "/home/brice", nil })
	if err != nil {
		t.Fatal(err)
	}
	if got != "/arg/cef" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveDirFallsBackToUserHome(t *testing.T) {
	t.Setenv("CEF_DIR", "")
	got, err := resolveDir("", func() (string, error) { return "/home/brice", nil })
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("/home/brice", ".local", "share", "cef")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestTargetMajorDefaultAndEnvOverride(t *testing.T) {
	os.Unsetenv("CEF_VERSION")
	if got := targetMajor(); got != 145 {
		t.Fatalf("default target = %d", got)
	}
	t.Setenv("CEF_VERSION", "146")
	if got := targetMajor(); got != 146 {
		t.Fatalf("override target = %d", got)
	}
}

func TestValidateVersionSkipsWhenRequested(t *testing.T) {
	t.Setenv("CEF_SKIP_VERSION_CHECK", "1")
	if err := validateVersion(func(int32) int32 { return 999 }); err != nil {
		t.Fatal(err)
	}
}

func TestValidateVersionReturnsMismatch(t *testing.T) {
	t.Setenv("CEF_SKIP_VERSION_CHECK", "")
	t.Setenv("CEF_VERSION", "145")
	err := validateVersion(func(int32) int32 { return 144 })
	if err == nil {
		t.Fatal("expected mismatch")
	}
}
