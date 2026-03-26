// internal/loader/loader.go
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ebitengine/purego"
)

const defaultCEFVersion = 145
const cefVersionInfoChromeMajor = 4

// Open loads libcef.so and validates the CEF version.
// The returned handle is used to register all C API symbols.
func Open(cefDir string) (uintptr, error) {
	runtimeDir, err := resolveDir(cefDir)
	if err != nil {
		return 0, err
	}
	libPath := filepath.Join(runtimeDir, "libcef.so")
	handle, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, fmt.Errorf("dlopen %s: %w", libPath, err)
	}
	if os.Getenv("CEF_SKIP_VERSION_CHECK") != "1" {
		if err := validateVersion(handle); err != nil {
			return 0, err
		}
	}
	return handle, nil
}

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

func validateVersion(handle uintptr) error {
	sym, err := purego.Dlsym(handle, "cef_version_info")
	if err != nil {
		return fmt.Errorf("resolve cef_version_info: %w", err)
	}
	var versionInfo func(int32) int32
	purego.RegisterFunc(&versionInfo, sym)
	got := versionInfo(cefVersionInfoChromeMajor)
	want := targetMajor()
	if got != want {
		return fmt.Errorf("unsupported CEF runtime: chrome major=%d want=%d", got, want)
	}
	return nil
}
