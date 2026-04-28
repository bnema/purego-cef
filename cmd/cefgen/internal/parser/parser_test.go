package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseStructWithCallbacks(t *testing.T) {
	cases := []string{
		`typedef struct _cef_client_t {
  cef_base_ref_counted_t base;
  struct _cef_render_handler_t*(CEF_CALLBACK* get_render_handler)(struct _cef_client_t* self);
} cef_client_t;`,
		`typedef struct _cef_accessibility_handler_t {
  cef_base_ref_counted_t base;
  void(CEF_CALLBACK *on_accessibility_tree_change)(struct _cef_accessibility_handler_t *self, struct _cef_value_t *value);
} cef_accessibility_handler_t;`,
	}

	for _, input := range cases {
		file, err := Parse("cef_client_capi.h", []byte(input))
		if err != nil {
			t.Fatal(err)
		}
		if len(file.Structs) != 1 || len(file.Structs[0].Fields) != 2 {
			t.Fatalf("unexpected parse result: %+v", file)
		}
		if !file.Structs[0].Fields[1].IsFunction {
			t.Fatalf("expected callback field, got: %+v", file.Structs[0].Fields[1])
		}
	}
}

func TestParseExportedFunctions(t *testing.T) {
	input := `CEF_EXPORT int cef_execute_process(const cef_main_args_t* args, cef_app_t* application, void* windows_sandbox_info);`
	file, err := Parse("cef_app_capi.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Functions) != 1 || file.Functions[0].CName != "cef_execute_process" {
		t.Fatalf("unexpected functions: %+v", file.Functions)
	}
}

func TestMapType(t *testing.T) {
	tests := []struct {
		ctype string
		want  string
	}{
		// Existing types (regression)
		{"int", "int32"},
		{"size_t", "uintptr"},
		{"void*", "unsafe.Pointer"},
		{"void", ""},
		{"double", "float64"},
		{"const char*", "*byte"},
		{"cef_base_ref_counted_t", "CEFBaseRefCountedT"},

		// New: cef_string_t by value -> CEFStringT (not uintptr)
		{"cef_string_t", "CEFStringT"},

		// Opaque string handles -> uintptr
		{"cef_string_userfree_t", "uintptr"},
		{"cef_string_list_t", "uintptr"},
		{"cef_string_map_t", "uintptr"},
		{"cef_string_multimap_t", "uintptr"},

		// New: stdint types
		{"uint32_t", "uint32"},
		{"uint64_t", "uint64"},
		{"int64_t", "int64"},
		{"int32_t", "int32"},
		{"uint16_t", "uint16"},
		{"int16_t", "int16"},

		// New: unsigned long (cef_window_handle_t on linux)
		{"unsigned long", "uint64"},

		// Pointer to cef_string_t -> unsafe.Pointer (handled by isPtr)
		{"cef_string_t*", "unsafe.Pointer"},
		{"const cef_string_t*", "unsafe.Pointer"},
	}
	for _, tt := range tests {
		t.Run(tt.ctype, func(t *testing.T) {
			got := mapType(tt.ctype)
			if got != tt.want {
				t.Errorf("mapType(%q) = %q, want %q", tt.ctype, got, tt.want)
			}
		})
	}
}

func TestStripCommentsRemovesPreprocessor(t *testing.T) {
	input := []byte(`#ifndef GUARD_H
#define GUARD_H
#include "foo.h"
#if CEF_API_ADDED(14100)
typedef struct _cef_foo_t {
  int bar;
} cef_foo_t;
#endif
#ifdef __cplusplus
extern "C" {
#endif
`)
	got := string(stripComments(input))
	if strings.Contains(got, "#if") || strings.Contains(got, "#endif") || strings.Contains(got, "#ifdef") {
		t.Errorf("stripComments did not remove preprocessor directives:\n%s", got)
	}
	if !strings.Contains(got, "typedef struct") {
		t.Error("stripComments removed non-preprocessor code")
	}
}

func TestStripCommentsRemovesBlockComments(t *testing.T) {
	input := []byte(`/* This is a block comment */
typedef struct _cef_foo_t {
  /* field comment */
  int bar;
} cef_foo_t;`)
	got := string(stripComments(input))
	if strings.Contains(got, "/*") || strings.Contains(got, "*/") {
		t.Errorf("stripComments did not remove block comments:\n%s", got)
	}
	if !strings.Contains(got, "int bar") {
		t.Error("stripComments removed non-comment code")
	}
}

func TestParseTypeHeader(t *testing.T) {
	input := []byte(`
#ifdef __cplusplus
extern "C" {
#endif

typedef struct _cef_point_t {
  int x;
  int y;
} cef_point_t;

typedef enum {
  LOGSEVERITY_DEFAULT,
  LOGSEVERITY_VERBOSE,
  LOGSEVERITY_DEBUG = LOGSEVERITY_VERBOSE,
  LOGSEVERITY_DISABLE = 99
} cef_log_severity_t;

#ifdef __cplusplus
}
#endif
`)
	file, err := Parse("cef_types_geometry.h", input)
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(file.Structs))
	}
	if file.Structs[0].GoName != "CEFPointT" {
		t.Errorf("struct name = %q, want CEFPointT", file.Structs[0].GoName)
	}
	if len(file.Structs[0].Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(file.Structs[0].Fields))
	}
	if len(file.Enums) != 1 {
		t.Fatalf("expected 1 enum, got %d", len(file.Enums))
	}
	if len(file.Enums[0].Values) != 4 {
		t.Errorf("expected 4 enum values, got %d", len(file.Enums[0].Values))
	}
}

func TestParseRealHeaders(t *testing.T) {
	cefDir := os.Getenv("CEF_DIR")
	if cefDir == "" {
		t.Skip("CEF_DIR not set")
	}

	patterns := []string{
		filepath.Join(cefDir, "include", "internal", "cef_types*.h"),
		filepath.Join(cefDir, "include", "capi", "*_capi.h"),
		filepath.Join(cefDir, "include", "capi", "views", "*_capi.h"),
		filepath.Join(cefDir, "include", "capi", "test", "*_capi.h"),
	}

	var count int
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}
		for _, path := range matches {
			if strings.HasSuffix(path, "_wrappers.h") {
				continue
			}
			t.Run(filepath.Base(path), func(t *testing.T) {
				_, err := ParseFile(path)
				if err != nil {
					t.Fatal(err)
				}
			})
			count++
		}
	}
	t.Logf("parsed %d headers successfully", count)
}
