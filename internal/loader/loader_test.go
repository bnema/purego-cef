package loader

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveDirOrder mutates package-level test hooks (systemCEFRuntimeDir,
// pathExists, and userHomeDir), so do not call t.Parallel() here; parallel
// execution would introduce data races and make future regressions harder to diagnose.
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

func TestDefaultPathExistsOnlyMissingReturnsFalse(t *testing.T) {
	existing := t.TempDir()
	missing := filepath.Join(existing, "missing")
	invalid := string([]byte{0})

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "existing", path: existing, want: true},
		{name: "missing", path: missing, want: false},
		{name: "other stat error", path: invalid, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultPathExists(tt.path); got != tt.want {
				t.Fatalf("defaultPathExists(%q)=%v want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestEnsureMinimumVersion(t *testing.T) {
	if err := ensureMinimumVersion(150, 150); err != nil {
		t.Fatalf("ensureMinimumVersion(equal): %v", err)
	}
	if err := ensureMinimumVersion(151, 150); err != nil {
		t.Fatalf("ensureMinimumVersion(newer): %v", err)
	}
	if err := ensureMinimumVersion(149, 150); err == nil {
		t.Fatalf("ensureMinimumVersion(older) expected error")
	} else if !strings.Contains(err.Error(), "minimum=150") {
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
