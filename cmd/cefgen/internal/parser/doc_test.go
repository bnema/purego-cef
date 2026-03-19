package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractStructDoc(t *testing.T) {
	input := `///
/// Implement this structure to provide handler implementations.
///
typedef struct _cef_client_t {
  cef_base_ref_counted_t base;
} cef_client_t;`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(header.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(header.Structs))
	}
	st := header.Structs[0]
	if st.Doc != "Implement this structure to provide handler implementations." {
		t.Errorf("Doc = %q", st.Doc)
	}
}

func TestDetectClientSideKind(t *testing.T) {
	input := `///
/// Implement this structure to provide handler implementations.
///
/// NOTE: This struct is allocated client-side.
///
typedef struct _cef_client_t {
  cef_base_ref_counted_t base;
} cef_client_t;`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	st := header.Structs[0]
	if st.Kind != "handler" {
		t.Errorf("Kind = %q, want %q", st.Kind, "handler")
	}
}

func TestDetectDLLSideKind(t *testing.T) {
	input := `///
/// Structure used to represent a browser.
///
/// NOTE: This struct is allocated DLL-side.
///
typedef struct _cef_browser_t {
  cef_base_ref_counted_t base;
} cef_browser_t;`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	st := header.Structs[0]
	if st.Kind != "object" {
		t.Errorf("Kind = %q, want %q", st.Kind, "object")
	}
}

func TestNoAllocationCommentIsData(t *testing.T) {
	input := `///
/// Structure representing a point.
///
typedef struct _cef_point_t {
  int x;
  int y;
} cef_point_t;`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	st := header.Structs[0]
	if st.Kind != "data" {
		t.Errorf("Kind = %q, want %q", st.Kind, "data")
	}
}

func TestExtractFieldDoc(t *testing.T) {
	input := `///
/// Implement this structure to provide handler implementations.
///
/// NOTE: This struct is allocated client-side.
///
typedef struct _cef_client_t {
  ///
  /// Base structure.
  ///
  cef_base_ref_counted_t base;

  ///
  /// Return the handler for audio rendering events.
  ///
  struct _cef_audio_handler_t*(CEF_CALLBACK* get_audio_handler)(
      struct _cef_client_t* self);
} cef_client_t;`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(header.Structs) != 1 {
		t.Fatalf("expected 1 struct, got %d", len(header.Structs))
	}
	st := header.Structs[0]

	// Find the get_audio_handler field.
	found := false
	for _, f := range st.Fields {
		if f.CName == "get_audio_handler" {
			found = true
			if f.Doc != "Return the handler for audio rendering events." {
				t.Errorf("field Doc = %q", f.Doc)
			}
			break
		}
	}
	if !found {
		t.Error("get_audio_handler field not found")
	}
}

func TestStructInterfaceName(t *testing.T) {
	input := `typedef struct _cef_browser_t {
  cef_base_ref_counted_t base;
} cef_browser_t;`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	st := header.Structs[0]
	if st.InterfaceName != "Browser" {
		t.Errorf("InterfaceName = %q, want %q", st.InterfaceName, "Browser")
	}
}

func TestFunctionDoc(t *testing.T) {
	input := `///
/// This function should be called from the application entry point function to
/// execute a secondary process.
///
CEF_EXPORT int cef_execute_process(const cef_main_args_t* args, cef_app_t* application, void* windows_sandbox_info);`

	header, err := Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(header.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d", len(header.Functions))
	}
	fn := header.Functions[0]
	want := "This function should be called from the application entry point function to execute a secondary process."
	if fn.Doc != want {
		t.Errorf("function Doc = %q, want %q", fn.Doc, want)
	}
}

func TestCleanDocStripsNOTE(t *testing.T) {
	lines := []string{
		"Implement this structure to provide handler implementations.",
		"",
		"NOTE: This struct is allocated client-side.",
		"",
	}
	got := cleanDoc(lines)
	want := "Implement this structure to provide handler implementations."
	if got != want {
		t.Errorf("cleanDoc = %q, want %q", got, want)
	}
}

func TestClassifyKind(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			"client-side",
			[]string{"NOTE: This struct is allocated client-side."},
			"handler",
		},
		{
			"DLL-side",
			[]string{"NOTE: This struct is allocated DLL-side."},
			"object",
		},
		{
			"no allocation",
			[]string{"Structure representing a point."},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyKind(tt.lines)
			if got != tt.want {
				t.Errorf("classifyKind = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseRealHeaderDocs(t *testing.T) {
	cefDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "cef")
	clientPath := filepath.Join(cefDir, "include", "capi", "cef_client_capi.h")
	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		t.Skip("CEF headers not found at", clientPath)
	}

	header, err := ParseFile(clientPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(header.Structs) == 0 {
		t.Fatal("no structs parsed")
	}
	st := header.Structs[0]

	// Verify struct has doc populated.
	if st.Doc == "" {
		t.Error("cef_client_t Doc is empty")
	}
	t.Logf("cef_client_t Doc: %q", st.Doc)

	// Verify Kind is handler (client-side).
	if st.Kind != "handler" {
		t.Errorf("cef_client_t Kind = %q, want %q", st.Kind, "handler")
	}

	// Verify InterfaceName.
	if st.InterfaceName != "Client" {
		t.Errorf("cef_client_t InterfaceName = %q, want %q", st.InterfaceName, "Client")
	}

	// Verify callback fields have docs.
	foundWithDoc := 0
	for _, f := range st.Fields {
		if f.IsFunction && f.Doc != "" {
			foundWithDoc++
		}
	}
	if foundWithDoc == 0 {
		t.Error("no callback fields have Doc populated")
	}
	t.Logf("%d/%d callback fields have docs", foundWithDoc, len(st.Fields)-1) // -1 for base

	// Check a specific field.
	for _, f := range st.Fields {
		if f.CName == "get_audio_handler" {
			if f.Doc == "" {
				t.Error("get_audio_handler Doc is empty")
			}
			t.Logf("get_audio_handler Doc: %q", f.Doc)
			break
		}
	}
}

func TestParseRealBrowserHeaderDocs(t *testing.T) {
	cefDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "cef")
	browserPath := filepath.Join(cefDir, "include", "capi", "cef_browser_capi.h")
	if _, err := os.Stat(browserPath); os.IsNotExist(err) {
		t.Skip("CEF headers not found at", browserPath)
	}

	header, err := ParseFile(browserPath)
	if err != nil {
		t.Fatal(err)
	}

	// cef_browser_t should be DLL-side -> object.
	for _, st := range header.Structs {
		if st.CName == "cef_browser_t" {
			if st.Kind != "object" {
				t.Errorf("cef_browser_t Kind = %q, want %q", st.Kind, "object")
			}
			if st.Doc == "" {
				t.Error("cef_browser_t Doc is empty")
			}
			t.Logf("cef_browser_t Kind=%q Doc=%q", st.Kind, st.Doc)
			break
		}
	}
}
