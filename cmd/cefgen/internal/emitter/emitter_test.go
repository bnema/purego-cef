package emitter

import (
	"strings"
	"testing"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
	"github.com/bnema/purego-cef/cmd/cefgen/internal/parser"
)

func TestEmitStructIncludesHostLayout(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:  "cef_client_t",
			GoName: "CEFClientT",
			Fields: []model.Field{{GoName: "Base", GoType: "CEFBaseRefCountedT"}},
		}},
	}
	code, err := Emit(header)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, "structs.HostLayout") {
		t.Fatal("missing HostLayout")
	}
}

func TestEmitRegisterFunction(t *testing.T) {
	header := &model.Header{Functions: []model.Function{{CName: "cef_initialize", GoName: "CEFInitialize"}}}
	code, err := Emit(header)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, `purego.RegisterLibFunc(&CEFInitialize, handle, "cef_initialize")`) {
		t.Fatal("missing RegisterLibFunc")
	}
}

func TestEmitFromParsedHeader(t *testing.T) {
	input := `typedef struct _cef_client_t {
  cef_base_ref_counted_t base;
  struct _cef_render_handler_t*(CEF_CALLBACK* get_render_handler)(struct _cef_client_t* self);
} cef_client_t;

CEF_EXPORT int cef_execute_process(const cef_main_args_t* args, cef_app_t* application, void* windows_sandbox_info);`

	header, err := parser.Parse("test.h", []byte(input))
	if err != nil {
		t.Fatal(err)
	}
	code, err := Emit(header)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, "CEFClientT") {
		t.Fatal("missing struct")
	}
	if !strings.Contains(code, "CEFExecuteProcess") {
		t.Fatal("missing function")
	}
	if !strings.Contains(code, "OverrideGetRenderHandler") {
		t.Fatal("missing Override method")
	}
	if !strings.Contains(code, "CallGetRenderHandler") {
		t.Fatal("missing Call method")
	}
}
