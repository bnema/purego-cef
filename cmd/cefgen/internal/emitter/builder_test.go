package emitter

import (
	"testing"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

func TestParamOverrideApplied(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_render_handler_t",
			GoName:        "CEFRenderHandlerT",
			Kind:          "handler",
			InterfaceName: "RenderHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "get_screen_point",
					GoName:      "GetScreenPoint",
					IsFunction:  true,
					ReturnCType: "int",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_render_handler_t*"},
						{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
						{CName: "screenX", GoName: "screenX", CType: "int*"},
						{CName: "screenY", GoName: "screenY", CType: "int*"},
					},
				},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	if len(data.Interfaces) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(data.Interfaces))
	}
	iface := data.Interfaces[0]
	if len(iface.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(iface.Methods))
	}
	m := iface.Methods[0]

	// Find screenX and screenY params by name.
	found := map[string]bool{"screenX": false, "screenY": false}
	for _, p := range m.Params {
		if p.Name == "screenX" || p.Name == "screenY" {
			found[p.Name] = true
			if p.PublicType != "*int32" {
				t.Errorf("param %s: expected PublicType *int32, got %s", p.Name, p.PublicType)
			}
			if p.MarshalKind != "dataStruct" {
				t.Errorf("param %s: expected MarshalKind dataStruct, got %s", p.Name, p.MarshalKind)
			}
		}
	}
	for name, ok := range found {
		if !ok {
			t.Errorf("param %s not found in method params", name)
		}
	}
}

func TestPixelBufferOverrideResolvesExpression(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_render_handler_t",
			GoName:        "CEFRenderHandlerT",
			Kind:          "handler",
			InterfaceName: "RenderHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "on_paint",
					GoName:      "OnPaint",
					IsFunction:  true,
					ReturnCType: "void",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_render_handler_t*"},
						{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
						{CName: "buffer", GoName: "buffer", CType: "const void*"},
						{CName: "width", GoName: "width", CType: "int"},
						{CName: "height", GoName: "height", CType: "int"},
					},
				},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	if len(data.Interfaces) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(data.Interfaces))
	}
	m := data.Interfaces[0].Methods[0]

	// Find buffer param.
	var bufParam *ParamData
	for i := range m.Params {
		if m.Params[i].Name == "buffer" {
			bufParam = &m.Params[i]
			break
		}
	}
	if bufParam == nil {
		t.Fatal("buffer param not found")
	}
	if bufParam.PublicType != "[]byte" {
		t.Errorf("expected PublicType []byte, got %s", bufParam.PublicType)
	}
	if bufParam.MarshalKind != "pixelBuffer" {
		t.Errorf("expected MarshalKind pixelBuffer, got %s", bufParam.MarshalKind)
	}
	if bufParam.UnmarshalExtra == "" {
		t.Error("expected non-empty UnmarshalExtra")
	}
	// The expression should have resolved {{width}} and {{height}} to argN names.
	// Params after skipping self: browser(arg0), buffer(arg1), width(arg2), height(arg3)
	want := "int(arg2)*int(arg3)*4"
	if bufParam.UnmarshalExtra != want {
		t.Errorf("expected UnmarshalExtra %q, got %q", want, bufParam.UnmarshalExtra)
	}
}
