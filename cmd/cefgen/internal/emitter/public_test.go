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
		{"interface declaration", "type Browser = portin.Browser"},
		{"method GetHost", "GetHost()"},
		{"method CanGoBack", "CanGoBack() bool"},
		{"impl struct", "type browserImpl struct"},
		{"wrap function", "func wrapBrowser(ptr unsafe.Pointer) Browser"},
		{"Release method", "func (obj *browserImpl) Release()"},
		{"runtime import", `"runtime"`},
		{"capi import", `capi`},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitPublicObjectInterface_UsesTypedPuregoForFloatMethods(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "_cef_browser_host_t",
			GoName:        "CEFBrowserHostT",
			Kind:          "object",
			InterfaceName: "BrowserHost",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "get_zoom_level",
					GoName:      "GetZoomLevel",
					IsFunction:  true,
					ReturnCType: "double",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_browser_host_t*"},
					},
				},
				{
					CName:       "set_zoom_level",
					GoName:      "SetZoomLevel",
					IsFunction:  true,
					ReturnCType: "void",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_browser_host_t*"},
						{CName: "zoomLevel", GoName: "zoomlevel", CType: "double"},
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

	checks := []string{
		`"github.com/bnema/purego"`,
		`var fn func(*capi.CEFBrowserHostT) float64`,
		`purego.RegisterFunc(&fn, obj.rawPtr.GetZoomLevel)`,
		`var fn func(*capi.CEFBrowserHostT, float64)`,
		`purego.RegisterFunc(&fn, obj.rawPtr.SetZoomLevel)`,
	}

	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}

	if strings.Contains(code, "math.Float64bits(zoomlevel)") {
		t.Fatalf("generated code still marshals float64 through math.Float64bits:\n%s", code)
	}
}

func TestEmitPublicHandlerInterface(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_focus_handler_t",
			GoName:        "CEFFocusHandlerT",
			Kind:          "handler",
			InterfaceName: "FocusHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "on_got_focus",
					GoName:      "OnGotFocus",
					IsFunction:  true,
					ReturnCType: "void",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_focus_handler_t*"},
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
		{"interface declaration", "type FocusHandler = portin.FocusHandler"},
		{"constructor function", "func NewFocusHandler(impl FocusHandler) FocusHandler"},
		{"initRefCount call", "initRefCount(unsafe.Pointer(r)"},
		{"purego.NewCallback", "purego.NewCallback"},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitPublicHandlerInterface_UsesTypedFloatCallbackParams(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_display_handler_t",
			GoName:        "CEFDisplayHandlerT",
			Kind:          "handler",
			InterfaceName: "DisplayHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "on_loading_progress_change",
					GoName:      "OnLoadingProgressChange",
					IsFunction:  true,
					ReturnCType: "void",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_display_handler_t*"},
						{CName: "progress", GoName: "progress", CType: "double"},
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

	if !strings.Contains(code, "purego.NewCallback(func(self uintptr, arg0 float64) {") {
		t.Fatalf("expected typed float callback parameter, got:\n%s", code)
	}

	if strings.Contains(code, "math.Float64frombits") {
		t.Fatalf("expected callback float parameter to avoid math.Float64frombits, got:\n%s", code)
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
		{"type declaration", "type State = int32"},
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
	if !strings.Contains(code, "capi.CEFGetMimeType") {
		t.Error("missing raw function reference")
	}
}

func TestEmitPublicDataStructUsesRenamedPublicType(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:  "cef_settings_t",
			GoName: "CEFSettingsT",
			Kind:   "data",
			Fields: []model.Field{{CName: "size", GoName: "Size", CType: "size_t", GoType: "uintptr"}},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	if !strings.Contains(code, "type CEFSettings = capi.CEFSettingsT") {
		t.Fatalf("expected renamed public data struct, got:\n%s", code)
	}
	if !strings.Contains(code, "func NewCEFSettings() CEFSettings") {
		t.Fatalf("expected constructor for renamed data struct, got:\n%s", code)
	}
}

func TestEmitPublicFreeFuncUsesNamedStringHandleTypes(t *testing.T) {
	header := &model.Header{
		Functions: []model.Function{{
			CName:       "cef_fill_string_list",
			GoName:      "CEFFillStringList",
			ReturnCType: "void",
			Params:      []model.Param{{CName: "values", GoName: "values", CType: "cef_string_list_t", GoType: "uintptr"}},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	if !strings.Contains(code, "func FillStringList(values StringList)") {
		t.Fatalf("expected named string-list handle in public signature, got:\n%s", code)
	}
	if !strings.Contains(code, "capi.CEFFillStringList(uintptr(values))") {
		t.Fatalf("expected raw call to cast named handle back to uintptr, got:\n%s", code)
	}
}

func TestEmitPublicDataStructUsesRenamedMainArgsType(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:  "cef_main_args_t",
			GoName: "CEFMainArgsT",
			Kind:   "data",
			Fields: []model.Field{{CName: "argc", GoName: "Argc", CType: "int", GoType: "int32"}},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	if !strings.Contains(code, "type CEFMainArgs = capi.CEFMainArgsT") {
		t.Fatalf("expected renamed public main args type, got:\n%s", code)
	}
}

func TestEmitPublicFreeFuncAutoWrapsHandlerParams(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_client_t",
			GoName:        "CEFClientT",
			Kind:          "handler",
			InterfaceName: "Client",
			Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
		}},
		Functions: []model.Function{{
			CName:       "cef_make_client_consumer",
			GoName:      "CEFMakeClientConsumer",
			ReturnCType: "void",
			Params:      []model.Param{{CName: "client", GoName: "client", CType: "struct _cef_client_t*"}},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	want := "extractOrWrapRawPointer(client, func() any { return newRawClient(client) })"
	if !strings.Contains(code, want) {
		t.Fatalf("expected generated code to contain %q\n\nGot:\n%s", want, code)
	}
	if strings.Contains(code, "extractRawPointer(client)") {
		t.Fatalf("expected generated code to auto-wrap client handler params, got:\n%s", code)
	}
}

func TestEmitPublicObjectMethodAutoWrapsHandlerParams(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "cef_task_t",
				GoName:        "CEFTaskT",
				Kind:          "handler",
				InterfaceName: "Task",
				Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
			},
			{
				CName:         "cef_task_runner_t",
				GoName:        "CEFTaskRunnerT",
				Kind:          "object",
				InterfaceName: "TaskRunner",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "post_task",
						GoName:      "PostTask",
						IsFunction:  true,
						ReturnCType: "int",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_task_runner_t*"},
							{CName: "task", GoName: "task", CType: "struct _cef_task_t*"},
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

	want := "CallPostTask(uintptr(extractOrWrapRawPointer(task, func() any { return NewTask(task) })))"
	if !strings.Contains(code, want) {
		t.Fatalf("expected generated code to contain %q\n\nGot:\n%s", want, code)
	}
	if strings.Contains(code, "CallPostTask(uintptr(extractRawPointer(task)))") {
		t.Fatalf("expected generated code to auto-wrap task handler params, got:\n%s", code)
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
		{"unsafe.Slice in unmarshal", "unsafe.Slice((*byte)(unsafe.Pointer("},
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

func TestEmitV8HandlerExecuteArgumentsObjectSlice(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "_cef_v8_value_t",
				GoName:        "CEFV8ValueT",
				Kind:          "object",
				InterfaceName: "V8Value",
				Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
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

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	checks := []string{
		"argumentsPtrs := unsafe.Slice((*uintptr)(unsafe.Pointer(arg3)), int(arg2))",
		"arguments[i] = wrapV8Value(unsafe.Pointer(ptr))",
		"return uintptr(impl.Execute(name, object, arguments, retval, exception))",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}

	if strings.Contains(code, "argumentsCount") || strings.Contains(code, "argumentscount") {
		t.Fatalf("expected argumentsCount to be merged away, got:\n%s", code)
	}
}

func TestEmitV8ContextEvalAddsNilOutputGuards(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "_cef_v8_context_t",
			GoName:        "CEFV8ContextT",
			Kind:          "object",
			InterfaceName: "V8Context",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "eval",
					GoName:      "Eval",
					IsFunction:  true,
					ReturnCType: "int",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_v8_context_t*"},
						{CName: "code", GoName: "code", CType: "const cef_string_t*"},
						{CName: "scriptURL", GoName: "scriptURL", CType: "const cef_string_t*"},
						{CName: "startLine", GoName: "startLine", CType: "int"},
						{CName: "retval", GoName: "retval", CType: "struct _cef_v8_value_t**"},
						{CName: "exception", GoName: "exception", CType: "struct _cef_v8_exception_t**"},
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

	checks := []string{
		"// CEF's eval requires valid output pointers",
		"if retval == nil {",
		"retval = unsafe.Pointer(&scratch)",
		"if exception == nil {",
		"exception = unsafe.Pointer(&scratch)",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}
}

func TestEmitV8ValueExecuteFunctionArgumentsObjectSlice(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "cef_v8_context_t",
				GoName:        "CEFV8ContextT",
				Kind:          "object",
				InterfaceName: "V8Context",
				Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
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

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	checks := []string{
		"func (obj *v8ValueImpl) ExecuteFunction(object V8Value, arguments []V8Value) V8Value {",
		"argumentsRaw = make([]uintptr, len(arguments))",
		"argumentsRaw[i] = uintptr(extractRawPointer(elem))",
		"argumentsPtr = unsafe.Pointer(&argumentsRaw[0])",
		"CallExecuteFunction(uintptr(extractRawPointer(object)), uintptr(len(arguments)), uintptr(argumentsPtr))",
		"func (obj *v8ValueImpl) ExecuteFunctionWithContext(context V8Context, object V8Value, arguments []V8Value) V8Value {",
		"CallExecuteFunctionWithContext(uintptr(extractRawPointer(context)), uintptr(extractRawPointer(object)), uintptr(len(arguments)), uintptr(argumentsPtr))",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}

	if strings.Contains(code, "argumentsCount") || strings.Contains(code, "argumentscount") {
		t.Fatalf("expected argumentsCount to be merged away, got:\n%s", code)
	}
}
