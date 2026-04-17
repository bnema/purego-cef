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

func TestObjectSliceMerge(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_request_handler_t",
			GoName:        "CEFRequestHandlerT",
			Kind:          "handler",
			InterfaceName: "RequestHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "on_select_client_certificate",
					GoName:      "OnSelectClientCertificate",
					IsFunction:  true,
					ReturnCType: "int",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_request_handler_t*"},
						{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
						{CName: "isProxy", GoName: "isProxy", CType: "int"},
						{CName: "host", GoName: "host", CType: "const cef_string_t*"},
						{CName: "port", GoName: "port", CType: "int"},
						{CName: "certificatesCount", GoName: "certificatesCount", CType: "size_t"},
						{CName: "certificates", GoName: "certificates", CType: "struct _cef_x509certificate_t* const*"},
						{CName: "callback", GoName: "callback", CType: "struct _cef_select_client_certificate_callback_t*"},
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

	// certificates param should have been merged with certificatesCount.
	var certParam *ParamData
	for i := range m.Params {
		if m.Params[i].Name == "certificates" {
			certParam = &m.Params[i]
		}
		if m.Params[i].Name == "certificatesCount" {
			t.Error("certificatesCount param should have been merged away, but is still present")
		}
	}
	if certParam == nil {
		t.Fatal("certificates param not found")
	}
	if certParam.PublicType != "[]X509Certificate" {
		t.Errorf("expected PublicType []X509Certificate, got %s", certParam.PublicType)
	}
	if certParam.MarshalKind != "objectSlice" {
		t.Errorf("expected MarshalKind objectSlice, got %s", certParam.MarshalKind)
	}
}

func TestV8HandlerExecuteArgumentsMerge(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "_cef_v8_value_t",
				GoName:        "CEFV8ValueT",
				Kind:          "object",
				InterfaceName: "V8Value",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				},
			},
			{
				CName:         "cef_v8_handler_t",
				GoName:        "CEFV8HandlerT",
				Kind:          "handler",
				InterfaceName: "V8Handler",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "execute",
						GoName:      "Execute",
						IsFunction:  true,
						ReturnCType: "int",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_v8_handler_t*"},
							{CName: "name", GoName: "name", CType: "const cef_string_t*"},
							{CName: "object", GoName: "object", CType: "struct _cef_v8_value_t*"},
							{CName: "argumentsCount", GoName: "argumentsCount", CType: "size_t"},
							{CName: "arguments", GoName: "arguments", CType: "struct _cef_v8_value_t* const*"},
							{CName: "retval", GoName: "retval", CType: "struct _cef_v8_value_t**"},
							{CName: "exception", GoName: "exception", CType: "cef_string_t*"},
						},
					},
				},
			},
		},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	if len(data.Interfaces) != 2 {
		t.Fatalf("expected 2 interfaces, got %d", len(data.Interfaces))
	}
	m := data.Interfaces[1].Methods[0]

	var argsParam *ParamData
	for i := range m.Params {
		if m.Params[i].Name == "arguments" {
			argsParam = &m.Params[i]
		}
		if m.Params[i].Name == "argumentsCount" {
			t.Error("argumentsCount param should have been merged away, but is still present")
		}
	}
	if argsParam == nil {
		t.Fatal("arguments param not found")
	}
	if argsParam.PublicType != "[]V8Value" {
		t.Errorf("expected PublicType []V8Value, got %s", argsParam.PublicType)
	}
	if argsParam.MarshalKind != "objectSlice" {
		t.Errorf("expected MarshalKind objectSlice, got %s", argsParam.MarshalKind)
	}
}

func TestV8ValueExecuteFunctionArgumentsMerge(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "cef_v8_context_t",
				GoName:        "CEFV8ContextT",
				Kind:          "object",
				InterfaceName: "V8Context",
				Fields: []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
			},
			{
				CName:         "cef_v8_value_t",
				GoName:        "CEFV8ValueT",
				Kind:          "object",
				InterfaceName: "V8Value",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "execute_function",
						GoName:      "ExecuteFunction",
						IsFunction:  true,
						ReturnCType: "struct _cef_v8_value_t*",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_v8_value_t*"},
							{CName: "object", GoName: "object", CType: "struct _cef_v8_value_t*"},
							{CName: "argumentsCount", GoName: "argumentsCount", CType: "size_t"},
							{CName: "arguments", GoName: "arguments", CType: "struct _cef_v8_value_t* const*"},
						},
					},
					{
						CName:       "execute_function_with_context",
						GoName:      "ExecuteFunctionWithContext",
						IsFunction:  true,
						ReturnCType: "struct _cef_v8_value_t*",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_v8_value_t*"},
							{CName: "context", GoName: "context", CType: "struct _cef_v8_context_t*"},
							{CName: "object", GoName: "object", CType: "struct _cef_v8_value_t*"},
							{CName: "argumentsCount", GoName: "argumentsCount", CType: "size_t"},
							{CName: "arguments", GoName: "arguments", CType: "struct _cef_v8_value_t* const*"},
						},
					},
				},
			},
		},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	var valueIface *InterfaceData
	for i := range data.Interfaces {
		if data.Interfaces[i].Name == "V8Value" {
			valueIface = &data.Interfaces[i]
			break
		}
	}
	if valueIface == nil {
		t.Fatal("V8Value interface not found")
	}
	if len(valueIface.Methods) != 2 {
		t.Fatalf("expected 2 V8Value methods, got %d", len(valueIface.Methods))
	}

	for _, m := range valueIface.Methods {
		var argsParam *ParamData
		for i := range m.Params {
			if m.Params[i].Name == "arguments" {
				argsParam = &m.Params[i]
			}
			if m.Params[i].Name == "argumentsCount" {
				t.Errorf("%s: argumentsCount param should have been merged away, but is still present", m.Name)
			}
		}
		if argsParam == nil {
			t.Fatalf("%s: arguments param not found", m.Name)
		}
		if argsParam.PublicType != "[]V8Value" {
			t.Errorf("%s: expected PublicType []V8Value, got %s", m.Name, argsParam.PublicType)
		}
		if argsParam.MarshalKind != "objectSlice" {
			t.Errorf("%s: expected MarshalKind objectSlice, got %s", m.Name, argsParam.MarshalKind)
		}
	}
}
