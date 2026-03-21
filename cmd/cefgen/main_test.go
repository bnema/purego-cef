package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegisterName(t *testing.T) {
	tests := []struct {
		base string
		want string
	}{
		{"cef_app_capi.h", "RegisterApp"},
		{"cef_browser_capi.h", "RegisterBrowser"},
		{"cef_life_span_handler_capi.h", "RegisterLifeSpanHandler"},
		{"cef_types.h", "RegisterTypes"},
		{"cef_types_geometry.h", "RegisterTypesGeometry"},
		{"cef_types_linux.h", "RegisterTypesLinux"},
	}
	for _, tt := range tests {
		t.Run(tt.base, func(t *testing.T) {
			got := registerName(tt.base)
			if got != tt.want {
				t.Errorf("registerName(%q) = %q, want %q", tt.base, got, tt.want)
			}
		})
	}
}

func TestOutputName(t *testing.T) {
	tests := []struct {
		relPath string
		want    string
	}{
		{"cef_browser_capi.h", "cef_browser.go"},
		{"views/cef_button_capi.h", "views_cef_button.go"},
		{"test/cef_test_helpers_capi.h", "test_cef_test_helpers.go"},
		{"cef_types.h", "cef_types.go"},
		{"cef_types_mac.h", "cef_types_darwin.go"},
		{"cef_types_win.h", "cef_types_windows.go"},
	}
	for _, tt := range tests {
		t.Run(tt.relPath, func(t *testing.T) {
			got := outputName(tt.relPath)
			if got != tt.want {
				t.Errorf("outputName(%q) = %q, want %q", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestPlatformRegisterName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"RegisterTypes", ""},
		{"RegisterTypesLinux", "linux"},
		{"RegisterTypesMac", "darwin"},
		{"RegisterTypesWin", "windows"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := platformRegisterName(tt.name)
			if got != tt.want {
				t.Errorf("platformRegisterName(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestWriteRawRegisterAggregator(t *testing.T) {
	dir := t.TempDir()
	names := []string{
		"RegisterTypes",
		"RegisterTypesLinux",
		"RegisterTypesMac",
		"RegisterTypesWin",
		"RegisterBase",
	}

	if err := writeRawRegisterAggregator(dir, names); err != nil {
		t.Fatalf("writeRawRegisterAggregator failed: %v", err)
	}

	checkFileContains := func(name string, wants ...string) {
		t.Helper()
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		content := string(data)
		for _, want := range wants {
			if !strings.Contains(content, want) {
				t.Fatalf("%s missing %q\n\n%s", name, want, content)
			}
		}
	}

	checkFileContains("register.go", "RegisterTypes(handle)", "RegisterBase(handle)", "registerPlatform(handle)")
	checkFileContains("register_platform_linux.go", "//go:build linux", "RegisterTypesLinux(handle)")
	checkFileContains("register_platform_darwin.go", "//go:build darwin", "RegisterTypesMac(handle)")
	checkFileContains("register_platform_windows.go", "//go:build windows", "RegisterTypesWin(handle)")
	checkFileContains("register_platform_default.go", "//go:build !linux && !darwin && !windows", "func registerPlatform(_ uintptr) {}")
}

func TestFilterOut(t *testing.T) {
	paths := []string{
		"/a/cef_types.h",
		"/a/cef_types_wrappers.h",
		"/a/cef_types_geometry.h",
	}
	got := filterOut(paths, "cef_types_wrappers.h")
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, p := range got {
		if p == "/a/cef_types_wrappers.h" {
			t.Error("filterOut did not remove cef_types_wrappers.h")
		}
	}
}

func TestDiscoverRegisterNamesFromRawDir(t *testing.T) {
	dir := t.TempDir()
	mustWriteFileContent(t, filepath.Join(dir, "cef_app.go"), "package raw\nfunc RegisterApp(_ uintptr) {}\n")
	mustWriteFileContent(t, filepath.Join(dir, "cef_browser.go"), "package raw\nfunc RegisterBrowser(_ uintptr) {}\n")
	mustWriteFileContent(t, filepath.Join(dir, "register.go"), "package raw\nfunc Register(_ uintptr) {}\n")
	mustWriteFileContent(t, filepath.Join(dir, "doc.go"), "package raw\n")

	got, err := discoverRegisterNamesFromRawDir(dir)
	if err != nil {
		t.Fatalf("discoverRegisterNamesFromRawDir failed: %v", err)
	}

	want := map[string]bool{
		"RegisterApp":     true,
		"RegisterBrowser": true,
		"Register":        false,
	}

	for name, shouldExist := range want {
		found := false
		for _, gotName := range got {
			if gotName == name {
				found = true
				break
			}
		}
		if found != shouldExist {
			t.Fatalf("register %s present=%v, want %v; got=%v", name, found, shouldExist, got)
		}
	}
}

func mustWriteFileContent(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := config{headersDir: "/tmp/headers", rawDir: "/tmp/raw", publicDir: "/tmp/pub", version: "145"}
	if err := cfg.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		cfg  config
	}{
		{"empty headersDir", config{headersDir: "", rawDir: "/tmp/raw", publicDir: "/tmp/pub"}},
		{"empty rawDir", config{headersDir: "/tmp/h", rawDir: "", publicDir: "/tmp/pub"}},
		{"empty publicDir", config{headersDir: "/tmp/h", rawDir: "/tmp/raw", publicDir: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.validate(); err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
