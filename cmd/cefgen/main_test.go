package main

import (
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

func TestConfigValidate(t *testing.T) {
	cfg := config{headersDir: "/tmp/headers", outputDir: "/tmp/out", version: "145"}
	if err := cfg.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		cfg  config
	}{
		{"empty headersDir", config{headersDir: "", outputDir: "/tmp/out", version: "145"}},
		{"empty outputDir", config{headersDir: "/tmp/headers", outputDir: "", version: "145"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.validate(); err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}
