package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ebitengine/purego"
)

const defaultCEFVersion = 145

// cefVersionInfoChromeMajor is the entry index for the chrome major version
// in cef_version_info().
const cefVersionInfoChromeMajor = 4

// Open loads the CEF shared library from the resolved directory and validates
// the runtime version. Currently Linux-only (libcef.so).
func Open(dir string) (uintptr, error) {
	runtimeDir, err := resolveDir(dir, os.UserHomeDir)
	if err != nil {
		return 0, err
	}
	libPath := filepath.Join(runtimeDir, "libcef.so")
	handle, err := purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, fmt.Errorf("dlopen %s: %w", libPath, err)
	}
	sym, err := purego.Dlsym(handle, "cef_version_info")
	if err != nil {
		return 0, fmt.Errorf("resolve cef_version_info: %w", err)
	}
	var versionInfo func(int32) int32
	purego.RegisterFunc(&versionInfo, sym)
	if err := validateVersion(versionInfo); err != nil {
		return 0, err
	}
	return handle, nil
}

func resolveDir(arg string, userHomeDir func() (string, error)) (string, error) {
	if env := os.Getenv("CEF_DIR"); env != "" {
		return env, nil
	}
	if arg != "" {
		return arg, nil
	}
	home, err := userHomeDir()
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

func validateVersion(versionInfo func(int32) int32) error {
	if os.Getenv("CEF_SKIP_VERSION_CHECK") == "1" {
		return nil
	}
	got := versionInfo(cefVersionInfoChromeMajor)
	want := targetMajor()
	if got != want {
		return fmt.Errorf("unsupported CEF runtime: chrome major=%d want=%d", got, want)
	}
	return nil
}
