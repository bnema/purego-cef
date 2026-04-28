package loader

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDirOrder(t *testing.T) {
	origCEFDir := getenv("CEF_DIR")
	t.Cleanup(func() {
		setenv("CEF_DIR", origCEFDir)
	})

	origSystemDir := systemCEFRuntimeDir
	origPathExists := pathExists
	origUserHomeDir := userHomeDir
	t.Cleanup(func() {
		systemCEFRuntimeDir = origSystemDir
		pathExists = origPathExists
		userHomeDir = origUserHomeDir
	})

	systemCEFRuntimeDir = "/system/cef"
	pathExists = func(string) bool { return true }
	userHomeDir = func() (string, error) { return "/home/test", nil }

	setenv("CEF_DIR", "/env/cef")
	got, err := resolveDir("/arg/cef")
	if err != nil {
		t.Fatalf("resolveDir(env): %v", err)
	}
	if got != "/env/cef" {
		t.Fatalf("resolveDir(env)=%q want %q", got, "/env/cef")
	}

	setenv("CEF_DIR", "")
	got, err = resolveDir("/arg/cef")
	if err != nil {
		t.Fatalf("resolveDir(arg): %v", err)
	}
	if got != "/arg/cef" {
		t.Fatalf("resolveDir(arg)=%q want %q", got, "/arg/cef")
	}

	got, err = resolveDir("")
	if err != nil {
		t.Fatalf("resolveDir(system): %v", err)
	}
	if got != "/system/cef" {
		t.Fatalf("resolveDir(system)=%q want %q", got, "/system/cef")
	}
}

func TestResolveDirFallbackHome(t *testing.T) {
	origCEFDir := getenv("CEF_DIR")
	t.Cleanup(func() {
		setenv("CEF_DIR", origCEFDir)
	})
	setenv("CEF_DIR", "")

	origSystemDir := systemCEFRuntimeDir
	origPathExists := pathExists
	origUserHomeDir := userHomeDir
	t.Cleanup(func() {
		systemCEFRuntimeDir = origSystemDir
		pathExists = origPathExists
		userHomeDir = origUserHomeDir
	})

	systemCEFRuntimeDir = "/system/cef"
	pathExists = func(string) bool { return false }
	userHomeDir = func() (string, error) { return "/home/test", nil }

	got, err := resolveDir("")
	if err != nil {
		t.Fatalf("resolveDir(fallback): %v", err)
	}
	want := filepath.Join("/home/test", ".local", "share", "cef")
	if got != want {
		t.Fatalf("resolveDir(fallback)=%q want %q", got, want)
	}
}

func TestResolveDirHomeError(t *testing.T) {
	origCEFDir := getenv("CEF_DIR")
	t.Cleanup(func() {
		setenv("CEF_DIR", origCEFDir)
	})
	setenv("CEF_DIR", "")

	origPathExists := pathExists
	origUserHomeDir := userHomeDir
	t.Cleanup(func() {
		pathExists = origPathExists
		userHomeDir = origUserHomeDir
	})

	pathExists = func(string) bool { return false }
	userHomeDir = func() (string, error) { return "", errors.New("no home") }

	if _, err := resolveDir(""); err == nil {
		t.Fatalf("resolveDir() expected error")
	}
}

func TestTargetMajorUsesMinimumOverride(t *testing.T) {
	orig := getenv("CEF_VERSION")
	t.Cleanup(func() { setenv("CEF_VERSION", orig) })

	setenv("CEF_VERSION", "149")
	if got := targetMajor(); got != 149 {
		t.Fatalf("targetMajor()=%d want 149", got)
	}

	setenv("CEF_VERSION", "not-a-number")
	if got := targetMajor(); got != defaultCEFVersion {
		t.Fatalf("targetMajor()=%d want default %d", got, defaultCEFVersion)
	}
}

func TestEnsureMinimumVersion(t *testing.T) {
	if err := ensureMinimumVersion(147, 147); err != nil {
		t.Fatalf("ensureMinimumVersion(equal): %v", err)
	}
	if err := ensureMinimumVersion(148, 147); err != nil {
		t.Fatalf("ensureMinimumVersion(newer): %v", err)
	}
	if err := ensureMinimumVersion(146, 147); err == nil {
		t.Fatalf("ensureMinimumVersion(older) expected error")
	} else if !strings.Contains(err.Error(), "minimum=147") {
		t.Fatalf("ensureMinimumVersion(older) error=%q want minimum text", err)
	}
}

func getenv(key string) string { return os.Getenv(key) }
func setenv(key, value string) {
	if value == "" {
		_ = os.Unsetenv(key)
		return
	}
	_ = os.Setenv(key, value)
}
