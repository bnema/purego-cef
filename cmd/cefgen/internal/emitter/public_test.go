package emitter

import (
	"strings"
	"testing"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

func TestEmitPublicObjectInterface(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "_cef_browser_t",
			GoName:        "CEFBrowserT",
			Kind:          "object",
			InterfaceName: "Browser",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "get_host",
					GoName:      "GetHost",
					IsFunction:  true,
					ReturnCType: "struct _cef_browser_host_t*",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_browser_t*"},
					},
				},
				{
					CName:       "can_go_back",
					GoName:      "CanGoBack",
					IsFunction:  true,
					ReturnCType: "int",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_browser_t*"},
					},
				},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	t.Log(code)

	checks := []struct {
		desc string
		want string
	}{
		{"interface declaration", "type Browser interface"},
		{"method GetHost", "GetHost()"},
		{"method CanGoBack", "CanGoBack() bool"},
		{"impl struct", "type browserImpl struct"},
		{"wrap function", "func wrapBrowser(ptr unsafe.Pointer) Browser"},
		{"Release method", "func (obj *browserImpl) Release()"},
		{"runtime import", `"runtime"`},
		{"raw import", `raw`},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitPublicHandlerInterface(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "_cef_life_span_handler_t",
			GoName:        "CEFLifeSpanHandlerT",
			Kind:          "handler",
			InterfaceName: "LifeSpanHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "on_before_popup",
					GoName:      "OnBeforePopup",
					IsFunction:  true,
					ReturnCType: "int",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_life_span_handler_t*"},
						{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
					},
				},
				{
					CName:       "on_after_created",
					GoName:      "OnAfterCreated",
					IsFunction:  true,
					ReturnCType: "void",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_life_span_handler_t*"},
						{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
					},
				},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	t.Log(code)

	checks := []struct {
		desc string
		want string
	}{
		{"interface declaration", "type LifeSpanHandler interface"},
		{"constructor function", "func NewLifeSpanHandler(impl LifeSpanHandler) LifeSpanHandler"},
		{"initRefCount call", "initRefCount(unsafe.Pointer(r)"},
		{"purego.NewCallback", "purego.NewCallback"},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitPublicEnum(t *testing.T) {
	header := &model.Header{
		Enums: []model.Enum{{
			CName:  "cef_state_t",
			GoName: "CEFStateT",
			Values: []model.EnumValue{
				{CName: "CEF_STATE_DEFAULT", GoName: "CEFStateDefault", Value: "0"},
				{CName: "CEF_STATE_ENABLED", GoName: "CEFStateEnabled", Value: "1"},
				{CName: "CEF_STATE_DISABLED", GoName: "CEFStateDisabled", Value: "2"},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	t.Log(code)

	checks := []struct {
		desc string
		want string
	}{
		{"type declaration", "type State int32"},
		{"default value", "StateDefault"},
		{"enabled value", "StateEnabled"},
		{"disabled value", "StateDisabled"},
		{"const block", "const ("},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitPublicFreeFunc(t *testing.T) {
	// Use a function name that's NOT in skipPublicTypes.
	header := &model.Header{
		Functions: []model.Function{{
			CName:       "cef_get_mime_type",
			GoName:      "CEFGetMimeType",
			ReturnCType: "cef_string_userfree_t",
			Params: []model.Param{
				{CName: "extension", GoName: "extension", CType: "const cef_string_t*"},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	t.Log(code)

	if !strings.Contains(code, "func GetMimeType(") {
		t.Error("missing free function wrapper")
	}
	if !strings.Contains(code, "raw.CEFGetMimeType") {
		t.Error("missing raw function reference")
	}
}

func TestBuildPublicFileData(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "_cef_browser_t",
				GoName:        "CEFBrowserT",
				Kind:          "object",
				InterfaceName: "Browser",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "get_host",
						GoName:      "GetHost",
						IsFunction:  true,
						ReturnCType: "struct _cef_browser_host_t*",
						Params:      []model.Param{{CName: "self", CType: "struct _cef_browser_t*"}},
					},
				},
			},
		},
		Enums: []model.Enum{
			{CName: "cef_state_t", GoName: "CEFStateT", Values: []model.EnumValue{
				{CName: "CEF_STATE_DEFAULT", Value: "0"},
			}},
		},
		Functions: []model.Function{
			// Use a non-skipped function name.
			{CName: "cef_get_mime_type", GoName: "CEFGetMimeType", ReturnCType: "cef_string_userfree_t"},
		},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	if len(data.Interfaces) != 1 {
		t.Fatalf("expected 1 interface, got %d", len(data.Interfaces))
	}
	if data.Interfaces[0].Name != "Browser" {
		t.Errorf("expected interface name Browser, got %s", data.Interfaces[0].Name)
	}
	if len(data.Interfaces[0].Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(data.Interfaces[0].Methods))
	}
	if data.Interfaces[0].Methods[0].Name != "GetHost" {
		t.Errorf("expected method name GetHost, got %s", data.Interfaces[0].Methods[0].Name)
	}

	if len(data.Enums) != 1 {
		t.Fatalf("expected 1 enum, got %d", len(data.Enums))
	}
	if data.Enums[0].Name != "State" {
		t.Errorf("expected enum name State, got %s", data.Enums[0].Name)
	}

	if len(data.FreeFunctions) != 1 {
		t.Fatalf("expected 1 free func, got %d", len(data.FreeFunctions))
	}
	if data.FreeFunctions[0].Name != "GetMimeType" {
		t.Errorf("expected func name GetMimeType, got %s", data.FreeFunctions[0].Name)
	}
}

func TestEmitPixelBufferOverride(t *testing.T) {
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
						{CName: "type", GoName: "type_", CType: "cef_paint_element_type_t"},
						{CName: "dirtyRectsCount", GoName: "dirtyRectsCount", CType: "size_t"},
						{CName: "dirtyRects", GoName: "dirtyRects", CType: "cef_rect_t const*"},
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

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	t.Log(code)

	checks := []struct {
		desc string
		want string
	}{
		{"buffer param is []byte", "buffer []byte"},
		{"unsafe.Slice in unmarshal", "unsafe.Slice"},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitObjectSliceOverride(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "_cef_x509certificate_t",
				GoName:        "CEFX509CertificateT",
				Kind:          "object",
				InterfaceName: "X509Certificate",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				},
			},
			{
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
			},
		},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	t.Log(code)

	checks := []struct {
		desc     string
		want     string
		mustHave bool
	}{
		{"certificates param is []X509Certificate", "certificates []X509Certificate", true},
		{"wrapX509Certificate in unmarshal", "wrapX509Certificate", true},
	}

	for _, c := range checks {
		if c.mustHave && !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}

	// Must NOT contain the count param in the interface signature.
	if strings.Contains(code, "certificatescount") || strings.Contains(code, "certificatesCount") {
		t.Error("count param should be merged away, but found certificatesCount in output")
	}
}
