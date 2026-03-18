package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseStructWithCallbacks(t *testing.T) {
	input := `typedef struct _cef_client_t {
  cef_base_ref_counted_t base;
  struct _cef_render_handler_t*(CEF_CALLBACK* get_render_handler)(struct _cef_client_t* self);
} cef_client_t;`

	file, err := Parse("cef_client_capi.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Structs) != 1 || len(file.Structs[0].Fields) != 2 {
		t.Fatalf("unexpected parse result: %+v", file)
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

func TestParseRealHeaders(t *testing.T) {
	cefDir := os.Getenv("CEF_DIR")
	if cefDir == "" {
		t.Skip("CEF_DIR not set")
	}
	headersDir := filepath.Join(cefDir, "include", "capi")
	for _, name := range []string{
		"cef_base_capi.h",
		"cef_app_capi.h",
		"cef_client_capi.h",
		"cef_render_handler_capi.h",
		"cef_life_span_handler_capi.h",
		"cef_browser_capi.h",
	} {
		if _, err := ParseFile(filepath.Join(headersDir, name)); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
}
