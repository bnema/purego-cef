package emitter

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

func TestMarshalCallArgsPreserveRawPointerTypesForBothConsumers(t *testing.T) {
	params := []ParamData{
		{Name: "iface", PublicType: "Browser", MarshalKind: "interface"},
		{Name: "opaque", PublicType: "unsafe.Pointer", MarshalKind: "interface"},
		{Name: "text", PublicType: "string", MarshalKind: "string"},
		{Name: "ptr", PublicType: "unsafe.Pointer", MarshalKind: "pointer"},
		{Name: "data", PublicType: "*Point", MarshalKind: "dataStruct"},
		{Name: "pixels", PublicType: "[]byte", MarshalKind: "pixelBuffer"},
		{Name: "itemscount", PublicType: "*int", MarshalKind: "outCount", CountPartnerName: "items"},
		{Name: "items", PublicType: "[]Point", MarshalKind: "outSlice", CountParamName: "itemscount"},
		{Name: "values", PublicType: "[]int32", MarshalKind: "slice"},
		{Name: "mode", PublicType: "Mode", MarshalKind: "enum"},
		{Name: "number", PublicType: "int32", MarshalKind: "numeric"},
		{Name: "opaqueHandle", PublicType: "uintptr", MarshalKind: "numeric", RawGoType: "unsafe.Pointer"},
	}
	const wantCall = "rawPtr.CallInvoke(extractRawPointer(iface), opaque, unsafe.Pointer(&textStr), ptr, unsafe.Pointer(data), pixelsPtr, unsafe.Pointer(itemsCountPtr), itemsPtr, uintptr(len(values)), valuesPtr, uintptr(mode), uintptr(number), unsafe.Pointer(opaqueHandle))"

	for _, kind := range []string{"object", "handler"} {
		t.Run(kind, func(t *testing.T) {
			method := MethodData{
				Name: "Invoke", RawFieldName: "Invoke", Params: params, RawParams: params,
				Return: ReturnData{IsVoid: true},
			}
			if kind == "handler" {
				// IsGetter suppresses the unrelated C-to-Go callback body so this focused
				// fixture can exercise the handler template's Go-to-C reverse wrapper.
				method.IsGetter = true
			}
			data := &PublicFileData{
				PackageName: "cef",
				Interfaces: []InterfaceData{{
					Name: "Consumer", Kind: kind, RawGoName: "CEFConsumerT", IsScoped: true,
					Methods: []MethodData{method},
				}},
			}

			code, err := EmitPublic(data)
			if err != nil {
				t.Fatalf("EmitPublic failed: %v", err)
			}
			if !strings.Contains(code, wantCall) {
				t.Fatalf("%s consumer did not preserve pointer arguments to raw wrapper; want %q\n\nGot:\n%s", kind, wantCall, code)
			}
		})
	}
}

func TestGeneratedRawAndPublicPointerConsumersCompileTogether(t *testing.T) {
	params := []ParamData{
		{Name: "iface", PublicType: "Consumer", MarshalKind: "interface"},
		{Name: "text", PublicType: "string", MarshalKind: "string"},
		{Name: "ptr", PublicType: "unsafe.Pointer", MarshalKind: "pointer"},
		{Name: "data", PublicType: "*Point", MarshalKind: "dataStruct"},
		{Name: "pixels", PublicType: "[]byte", MarshalKind: "pixelBuffer"},
		{Name: "itemscount", PublicType: "*int", MarshalKind: "outCount", CountPartnerName: "items"},
		{Name: "items", PublicType: "[]Point", MarshalKind: "outSlice", CountParamName: "itemscount"},
		{Name: "values", PublicType: "[]int32", MarshalKind: "slice", SliceElemType: "int32"},
		{Name: "mode", PublicType: "Mode", MarshalKind: "enum"},
		{Name: "number", PublicType: "int32", MarshalKind: "numeric"},
	}
	data := &PublicFileData{
		PackageName: "cef",
		Interfaces: []InterfaceData{{
			Name: "Consumer", Kind: "object", RawGoName: "CEFConsumerT", IsScoped: true,
			Methods: []MethodData{{
				Name: "Invoke", RawFieldName: "Invoke", Params: params, RawParams: params,
				Return: ReturnData{IsVoid: true},
			}},
		}},
		DataStructs: []DataStructData{{Name: "Point", RawGoName: "CEFPointT"}},
		Enums:       []EnumData{{Name: "Mode", RawGoName: "CEFModeT"}},
	}
	rawHeader := &model.Header{
		Path: "tiny_capi.h", Package: "raw", RegisterName: "RegisterTiny",
		Enums: []model.Enum{{GoName: "CEFModeT"}},
		Structs: []model.Struct{
			{GoName: "CEFPointT", Fields: []model.Field{{GoName: "X", GoType: "int32"}}},
			{GoName: "CEFConsumerT", Fields: []model.Field{{
				GoName: "Invoke", GoType: "uintptr", IsFunction: true,
				Params: []model.Param{
					{CName: "self", GoName: "self", CType: "cef_consumer_t*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "iface", GoName: "iface", CType: "cef_consumer_t*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "text", GoName: "text", CType: "cef_string_t*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "ptr", GoName: "ptr", CType: "void*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "data", GoName: "data", CType: "cef_point_t*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "pixels", GoName: "pixels", CType: "void*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "itemscount", GoName: "itemscount", CType: "size_t*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "items", GoName: "items", CType: "cef_point_t*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "valuescount", GoName: "valuescount", CType: "size_t", GoType: "uintptr"},
					{CName: "values", GoName: "values", CType: "int*", GoType: "unsafe.Pointer", Pointer: 1},
					{CName: "mode", GoName: "mode", CType: "cef_mode_t", GoType: "CEFModeT"},
					{CName: "number", GoName: "number", CType: "int", GoType: "int32"},
				},
			}}},
		},
	}

	rawCode, err := EmitRaw(rawHeader)
	if err != nil {
		t.Fatalf("EmitRaw failed: %v", err)
	}
	publicCode, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}
	portCode, err := EmitPortIn(data)
	if err != nil {
		t.Fatalf("EmitPortIn failed: %v", err)
	}

	root := t.TempDir()
	writeGeneratedTestFile(t, root, "go.mod", "module github.com/bnema/purego-cef\n\ngo 1.26\n\nrequire github.com/bnema/purego v0.11.0-bnema.4\n")
	writeGeneratedTestFile(t, root, "internal/capi/tiny_gen.go", rawCode)
	writeGeneratedTestFile(t, root, "internal/ports/in/tiny_gen.go", portCode)
	writeGeneratedTestFile(t, root, "cef/tiny_gen.go", publicCode)
	writeGeneratedTestFile(t, root, "cef/helpers_testbuild.go", `package cef

import "unsafe"

type tinyCEFString struct{}
func cefString(string) tinyCEFString { return tinyCEFString{} }
func freeCefString(*tinyCEFString) {}
func extractRawPointer(any) unsafe.Pointer { return nil }
`)

	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated raw and public packages do not compile: %v\n%s", err, output)
	}
}

func writeGeneratedTestFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

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
		`getZoomLevelFunc func(*capi.CEFBrowserHostT) float64`,
		`registerTypedCallback(&obj.getZoomLevelFunc, rawPtr.GetZoomLevel)`,
		`setZoomLevelFunc func(*capi.CEFBrowserHostT, float64)`,
		`registerTypedCallback(&obj.setZoomLevelFunc, rawPtr.SetZoomLevel)`,
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

func TestEmitPublicObjectInterface_GetClientReturnsWrappedRawClient(t *testing.T) {
	headers := []*model.Header{
		{
			Structs: []model.Struct{{
				CName:         "_cef_browser_host_t",
				GoName:        "CEFBrowserHostT",
				Kind:          "object",
				InterfaceName: "BrowserHost",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "get_client",
						GoName:      "GetClient",
						IsFunction:  true,
						ReturnCType: "struct _cef_client_t*",
						Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_browser_host_t*"}},
					},
				},
			}},
		},
		{
			Structs: []model.Struct{{
				CName:         "_cef_client_t",
				GoName:        "CEFClientT",
				Kind:          "handler",
				InterfaceName: "RawClient",
				Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
			}},
		},
	}

	registry := NewTypeRegistry(headers)
	data := BuildPublicFileData(headers[0], registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	if !strings.Contains(code, "func (obj *browserHostImpl) GetClient() RawClient") {
		t.Fatalf("GetClient should return RawClient:\n%s", code)
	}
	if !strings.Contains(code, "return wrapRawClient(unsafe.Pointer(ret))") {
		t.Fatalf("GetClient should wrap the raw client pointer:\n%s", code)
	}
}

func TestEmitPublicObjectInterface_ReleaseIsIdempotent(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "_cef_page_t",
			GoName:        "CEFPageT",
			Kind:          "object",
			InterfaceName: "Page",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "get_id",
					GoName:      "GetID",
					IsFunction:  true,
					ReturnCType: "int",
					Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_page_t*"}},
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
		"if obj.rawPtr == nil {",
		"runtime.SetFinalizer(obj, nil)",
		"runtime.SetFinalizer(impl, (*pageImpl).Release)",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Fatalf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
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
		{"initRefCount pins wrapper owner", "initRefCount(unsafe.Pointer(r), unsafe.Sizeof(*r), w)"},
		{"shared callback", "sharedCEFCallback(&focusHandlerOnGotFocusSharedOnce"},
		{"owner dispatch", "cefCallbackOwnerAs[FocusHandler](self)"},
		{"idempotent release guard", "if obj.rawPtr == nil {"},
		{"release clears finalizer", "runtime.SetFinalizer(obj, nil)"},
		{"release finalizer", "runtime.SetFinalizer(impl, (*focusHandlerImpl).Release)"},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}
}

func TestEmitPublicHandlerGetterUsesSharedCallbackAndCachedWrapperField(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "cef_client_t",
				GoName:        "CEFClientT",
				Kind:          "handler",
				InterfaceName: "RawClient",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "get_life_span_handler",
						GoName:      "GetLifeSpanHandler",
						IsFunction:  true,
						ReturnCType: "struct _cef_life_span_handler_t*",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_client_t*"},
						},
					},
				},
			},
			{
				CName:         "cef_life_span_handler_t",
				GoName:        "CEFLifeSpanHandlerT",
				Kind:          "handler",
				InterfaceName: "RawLifeSpanHandler",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
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

	checks := []struct {
		desc string
		want string
	}{
		{"cached wrapper field", "cachedGetLifeSpanHandlerPtr unsafe.Pointer"},
		{"shared getter callback", "sharedCEFCallback(&rawClientGetLifeSpanHandlerSharedOnce"},
		{"wrapper owner dispatch", "cefCallbackOwnerAs[*rawClientWrapper](self)"},
		{"getter addRef", "addRef(owner.cachedGetLifeSpanHandlerPtr)"},
		{"constructor stores cached pointer", "w.cachedGetLifeSpanHandlerPtr = extractOrWrapRawPointer"},
		{"constructor installs shared callback", "r.OverrideGetLifeSpanHandler(rawClientGetLifeSpanHandlerCEFCallback())"},
	}
	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}

	forbidden := "r.OverrideGetLifeSpanHandler(newCEFCallback"
	if strings.Contains(code, forbidden) {
		t.Errorf("getter callback still uses per-object newCEFCallback: found %q", forbidden)
	}
}

func TestEmitPublicRawAudioHandlerEmitsReverseWrapper(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{}}
	audio := model.Struct{
		CName:         "cef_audio_handler_t",
		GoName:        "CEFAudioHandlerT",
		Kind:          "handler",
		InterfaceName: "AudioHandler",
		Fields: []model.Field{
			{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
			{
				CName:       "get_audio_parameters",
				GoName:      "GetAudioParameters",
				IsFunction:  true,
				ReturnCType: "int",
				Params: []model.Param{
					{CName: "self", GoName: "self", CType: "struct _cef_audio_handler_t*"},
					{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
					{CName: "params", GoName: "params", CType: "struct _cef_audio_parameters_t*"},
				},
			},
		},
	}
	browser := model.Struct{CName: "_cef_browser_t", GoName: "CEFBrowserT", Kind: "object", InterfaceName: "Browser"}
	params := model.Struct{CName: "_cef_audio_parameters_t", GoName: "CEFAudioParametersT", Kind: "data", InterfaceName: "AudioParameters"}
	header.Structs = append(header.Structs, audio, browser, params)

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	checks := []struct {
		desc string
		want string
	}{
		{"raw alias", "type RawAudioHandler = portin.RawAudioHandler"},
		{"raw constructor", "func NewRawAudioHandler(impl RawAudioHandler) RawAudioHandler"},
		{"reverse wrapper name", "func wrapAudioHandler(ptr unsafe.Pointer) RawAudioHandler"},
		{"reverse wrapper delegates", "return int32(ret)"},
	}
	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output:\n%s", c.desc, c.want, code)
		}
	}
	if strings.Contains(code, "func wrapRawAudioHandler") {
		t.Fatalf("expected custom wrapAudioHandler name, got generated raw-name wrapper:\n%s", code)
	}
}

func TestEmitPublicRawLifeSpanHandlerEmitsReverseWrapper(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{}}
	lifeSpan := model.Struct{
		CName:         "cef_life_span_handler_t",
		GoName:        "CEFLifeSpanHandlerT",
		Kind:          "handler",
		InterfaceName: "LifeSpanHandler",
		Fields: []model.Field{
			{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
			{
				CName:       "do_close",
				GoName:      "DoClose",
				IsFunction:  true,
				ReturnCType: "int",
				Params: []model.Param{
					{CName: "self", GoName: "self", CType: "struct _cef_life_span_handler_t*"},
					{CName: "browser", GoName: "browser", CType: "struct _cef_browser_t*"},
				},
			},
		},
	}
	browser := model.Struct{CName: "_cef_browser_t", GoName: "CEFBrowserT", Kind: "object", InterfaceName: "Browser"}
	header.Structs = append(header.Structs, lifeSpan, browser)

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	checks := []struct {
		desc string
		want string
	}{
		{"raw alias", "type RawLifeSpanHandler = portin.RawLifeSpanHandler"},
		{"raw constructor", "func NewRawLifeSpanHandler(impl RawLifeSpanHandler) RawLifeSpanHandler"},
		{"reverse wrapper name", "func wrapLifeSpanHandler(ptr unsafe.Pointer) RawLifeSpanHandler"},
		{"bool reverse wrapper return", "return ret != 0"},
	}
	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output:\n%s", c.desc, c.want, code)
		}
	}
	if strings.Contains(code, "func wrapRawLifeSpanHandler") {
		t.Fatalf("expected custom wrapLifeSpanHandler name, got generated raw-name wrapper:\n%s", code)
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

	if !strings.Contains(code, "func(self uintptr, arg0 float64) {") {
		t.Fatalf("expected typed float callback parameter, got:\n%s", code)
	}

	if strings.Contains(code, "math.Float64frombits") {
		t.Fatalf("expected callback float parameter to avoid math.Float64frombits, got:\n%s", code)
	}
}

func TestEmitPublicHandlerInterface_UsesTypedInt64CallbackParams(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_audio_stream_handler_t",
			GoName:        "CEFAudioStreamHandlerT",
			Kind:          "handler",
			InterfaceName: "AudioStreamHandler",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "on_audio_stream_packet",
					GoName:      "OnAudioStreamPacket",
					IsFunction:  true,
					ReturnCType: "void",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_audio_stream_handler_t*"},
						{CName: "pts", GoName: "pts", CType: "int64_t"},
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
		"registerTypedCallback(&obj.onAudioStreamPacketFunc, rawPtr.",
		"func(self uintptr, arg0 int64)",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Fatalf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}

	if strings.Contains(code, "uintptr(pts)") {
		t.Fatalf("expected int64 callback arg to avoid uintptr(pts), got:\n%s", code)
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

	if !strings.Contains(code, "type RawSettings = capi.CEFSettingsT") {
		t.Fatalf("expected renamed public data struct, got:\n%s", code)
	}
	if !strings.Contains(code, "func NewRawSettings() RawSettings") {
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

	if !strings.Contains(code, "type RawMainArgs = capi.CEFMainArgsT") {
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

	want := "extractOrWrapRawPointer(client, func() any { return NewRawClient(client) })"
	if !strings.Contains(code, want) {
		t.Fatalf("expected generated code to contain %q\n\nGot:\n%s", want, code)
	}
	if strings.Contains(code, "extractRawPointer(client)") {
		t.Fatalf("expected generated code to auto-wrap client handler params, got:\n%s", code)
	}
}

func TestEmitGeneratedClientUsesRawAudioConstructor(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "cef_audio_handler_t",
				GoName:        "CEFAudioHandlerT",
				Kind:          "handler",
				InterfaceName: "AudioHandler",
				Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
			},
			{
				CName:         "cef_client_t",
				GoName:        "CEFClientT",
				Kind:          "handler",
				InterfaceName: "Client",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "get_audio_handler",
						GoName:      "GetAudioHandler",
						IsFunction:  true,
						ReturnCType: "struct _cef_audio_handler_t*",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_client_t*"},
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

	want := "return NewRawAudioHandler(h)"
	if !strings.Contains(code, want) {
		t.Fatalf("expected generated code to contain %q\n\nGot:\n%s", want, code)
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

	want := "CallPostTask(extractOrWrapRawPointer(task, func() any { return NewTask(task) }))"
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
		{"bufferPtr declaration", "var bufferPtr unsafe.Pointer"},
		{"bufferPtr nil/empty guard", "if len(buffer) > 0 {"},
		{"bufferPtr first element address", "bufferPtr = unsafe.Pointer(&buffer[0])"},
		{"bufferPtr preserved as pointer", "rawPtr.CallOnPaint(uintptr(browser), type_, uintptr(dirtyRectsCount), dirtyRects, bufferPtr"},
	}

	for _, c := range checks {
		if !strings.Contains(code, c.want) {
			t.Errorf("missing %s: want %q in output", c.desc, c.want)
		}
	}

	if strings.Contains(code, "unsafe.SliceData(buffer)") {
		t.Fatalf("expected CallOnPaint to avoid unsafe.SliceData(buffer), got:\n%s", code)
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
		// uintptr callback arguments may represent Go-backed test registrars. Keep
		// their conversion in the runtime's narrowly scoped nocheckptr bridge.
		"object := wrapV8Value(cefCallbackPointer(arg1))",
		"argumentsPtrs := unsafe.Slice((*uintptr)(cefCallbackPointer(arg3)), int(arg2))",
		"arguments[i] = wrapV8Value(cefCallbackPointer(ptr))",
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
		"CallExecuteFunction(extractRawPointer(object), uintptr(len(arguments)), argumentsPtr)",
		"func (obj *v8ValueImpl) ExecuteFunctionWithContext(context V8Context, object V8Value, arguments []V8Value) V8Value {",
		"CallExecuteFunctionWithContext(extractRawPointer(context), extractRawPointer(object), uintptr(len(arguments)), argumentsPtr)",
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

func TestEmitObjectMethodOutputObjectSliceBuffer(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:         "cef_post_data_element_t",
				GoName:        "CEFPostDataElementT",
				Kind:          "object",
				InterfaceName: "PostDataElement",
				Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
			},
			{
				CName:         "cef_post_data_t",
				GoName:        "CEFPostDataT",
				Kind:          "object",
				InterfaceName: "PostData",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
					{
						CName:       "get_elements",
						GoName:      "GetElements",
						IsFunction:  true,
						ReturnCType: "void",
						Params: []model.Param{
							{CName: "self", GoName: "self", CType: "struct _cef_post_data_t*"},
							{CName: "elementsCount", GoName: "elementsCount", CType: "size_t*"},
							{CName: "elements", GoName: "elements", CType: "struct _cef_post_data_element_t**"},
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
		"func (obj *postDataImpl) GetElements(elementsCount *int, elements []PostDataElement) {",
		"elementsCountPtr := elementsCount",
		"elementsRaw = make([]uintptr, len(elements))",
		"rawPtr.CallGetElements(unsafe.Pointer(elementsCountPtr), elementsPtr)",
		"elements[i] = wrapPostDataElement(unsafe.Pointer(elementsRaw[i]))",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}
}

func TestEmitObjectMethodOutputNumericSliceBuffer(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_task_manager_t",
			GoName:        "CEFTaskManagerT",
			Kind:          "object",
			InterfaceName: "TaskManager",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false},
				{
					CName:       "get_task_ids_list",
					GoName:      "GetTaskIdsList",
					IsFunction:  true,
					ReturnCType: "int",
					Params: []model.Param{
						{CName: "self", GoName: "self", CType: "struct _cef_task_manager_t*"},
						{CName: "task_idsCount", GoName: "task_idsCount", CType: "size_t*"},
						{CName: "task_ids", GoName: "task_ids", CType: "int64_t*"},
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
		"func (obj *taskManagerImpl) GetTaskIdsList(task_idsCount *int, task_ids []int64) int32 {",
		"task_idsCountPtr := task_idsCount",
		"task_idsPtr = unsafe.Pointer(&task_ids[0])",
		"ret := rawPtr.CallGetTaskIdsList(unsafe.Pointer(task_idsCountPtr), task_idsPtr)",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}
}

func TestEmitFreeFuncOutputObjectSliceBuffer(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:         "cef_display_t",
			GoName:        "CEFDisplayT",
			Kind:          "object",
			InterfaceName: "Display",
			Fields:        []model.Field{{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t", IsFunction: false}},
		}},
		Functions: []model.Function{{
			CName:       "cef_display_get_alls",
			GoName:      "CEFDisplayGetAlls",
			ReturnCType: "void",
			Params: []model.Param{
				{CName: "displaysCount", GoName: "displaysCount", CType: "size_t*", GoType: "unsafe.Pointer"},
				{CName: "displays", GoName: "displays", CType: "cef_display_t**", GoType: "unsafe.Pointer"},
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
		"func DisplayGetAlls(displaysCount *int, displays []Display) {",
		"displaysCountPtr := displaysCount",
		"displaysRaw = make([]uintptr, len(displays))",
		"capi.CEFDisplayGetAlls(unsafe.Pointer(displaysCountPtr), displaysPtr)",
		"displays[i] = wrapDisplay(unsafe.Pointer(displaysRaw[i]))",
	}
	for _, want := range checks {
		if !strings.Contains(code, want) {
			t.Errorf("expected generated code to contain %q\n\nGot:\n%s", want, code)
		}
	}
}

// TestEmitPublicGlobalFactoryUsesTakeOwnership verifies that a ref-counted
// interface returned by a global (free) function is adopted with takeX (no
// AddRef) rather than wrapX, and that a takeX constructor is emitted for it.
func TestEmitPublicGlobalFactoryUsesTakeOwnership(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:  "_cef_response_t",
			GoName: "CEFResponseT",
			Kind:   "object",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t"},
				{
					CName:       "is_read_only",
					GoName:      "IsReadOnly",
					IsFunction:  true,
					ReturnCType: "int",
					Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_response_t*"}},
				},
			},
		}},
		Functions: []model.Function{{
			CName:       "cef_response_create",
			GoName:      "CEFResponseCreate",
			ReturnCType: "struct _cef_response_t*",
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	if !strings.Contains(code, "func takeResponse(ptr unsafe.Pointer) Response {") {
		t.Errorf("expected takeResponse constructor to be emitted\n\nGot:\n%s", code)
	}
	// takeResponse must adopt without AddRef.
	takeBody := code[strings.Index(code, "func takeResponse"):]
	takeBody = takeBody[:strings.Index(takeBody, "\n}")]
	if strings.Contains(takeBody, "CallAddRef") {
		t.Errorf("takeResponse must NOT call AddRef\n\nGot:\n%s", takeBody)
	}
	if !strings.Contains(takeBody, "runtime.SetFinalizer(impl, (*responseImpl).Release)") {
		t.Errorf("takeResponse must set the Release finalizer\n\nGot:\n%s", takeBody)
	}
	// The global factory function must use takeResponse, not wrapResponse.
	if !strings.Contains(code, "func ResponseCreate() Response {") {
		t.Fatalf("missing ResponseCreate\n\nGot:\n%s", code)
	}
	createBody := code[strings.Index(code, "func ResponseCreate"):]
	createBody = createBody[:strings.Index(createBody, "\n}")]
	if !strings.Contains(createBody, "return takeResponse(ret)") {
		t.Errorf("ResponseCreate must adopt ownership via takeResponse\n\nGot:\n%s", createBody)
	}
	if strings.Contains(createBody, "wrapResponse") {
		t.Errorf("ResponseCreate must not use wrapResponse\n\nGot:\n%s", createBody)
	}
	// wrapResponse must remain unchanged and still AddRef borrowed pointers.
	wrapBody := code[strings.Index(code, "func wrapResponse"):]
	wrapBody = wrapBody[:strings.Index(wrapBody, "\n}")]
	if !strings.Contains(wrapBody, "base.CallAddRef()") {
		t.Errorf("wrapResponse must keep its AddRef\n\nGot:\n%s", wrapBody)
	}
}

// TestEmitPublicMethodReturnStillUsesWrap verifies that method returns of a
// ref-counted interface keep using wrapX even when a global factory for the same
// type exists (method-return ownership is explicitly out of scope for takeX).
func TestEmitPublicMethodReturnStillUsesWrap(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{
			{
				CName:  "_cef_response_t",
				GoName: "CEFResponseT",
				Kind:   "object",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t"},
					{
						CName:       "is_read_only",
						GoName:      "IsReadOnly",
						IsFunction:  true,
						ReturnCType: "int",
						Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_response_t*"}},
					},
				},
			},
			{
				CName:  "_cef_urlrequest_t",
				GoName: "CEFUrlrequestT",
				Kind:   "object",
				Fields: []model.Field{
					{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t"},
					{
						CName:       "get_response",
						GoName:      "GetResponse",
						IsFunction:  true,
						ReturnCType: "struct _cef_response_t*",
						Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_urlrequest_t*"}},
					},
				},
			},
		},
		Functions: []model.Function{{
			CName:       "cef_response_create",
			GoName:      "CEFResponseCreate",
			ReturnCType: "struct _cef_response_t*",
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	getResp := code[strings.Index(code, "func (obj *urlrequestImpl) GetResponse"):]
	getResp = getResp[:strings.Index(getResp, "\n}")]
	if !strings.Contains(getResp, "wrapResponse(unsafe.Pointer(ret))") {
		t.Errorf("method return GetResponse must keep wrapResponse\n\nGot:\n%s", getResp)
	}
	if strings.Contains(getResp, "takeResponse") {
		t.Errorf("method return GetResponse must NOT use takeResponse\n\nGot:\n%s", getResp)
	}
}

// TestEmitPublicNegativeSizeGuard verifies that a global function converting a
// signed Go int size param to a C size_t rejects negative values before the call.
func TestEmitPublicNegativeSizeGuard(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:  "_cef_binary_value_t",
			GoName: "CEFBinaryValueT",
			Kind:   "object",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t"},
				{
					CName:       "is_valid",
					GoName:      "IsValid",
					IsFunction:  true,
					ReturnCType: "int",
					Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_binary_value_t*"}},
				},
			},
		}},
		Functions: []model.Function{{
			CName:       "cef_binary_value_create",
			GoName:      "CEFBinaryValueCreate",
			ReturnCType: "struct _cef_binary_value_t*",
			Params: []model.Param{
				{CName: "data", GoName: "data", CType: "const void*", GoType: "unsafe.Pointer"},
				{CName: "data_size", GoName: "dataSize", CType: "size_t", GoType: "uintptr"},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	data := BuildPublicFileData(header, registry)

	code, err := EmitPublic(data)
	if err != nil {
		t.Fatalf("EmitPublic failed: %v", err)
	}

	createBody := code[strings.Index(code, "func BinaryValueCreate"):]
	createBody = createBody[:strings.Index(createBody, "\n}")]
	if !strings.Contains(createBody, "if dataSize < 0 {") {
		t.Errorf("expected negative-size guard for dataSize\n\nGot:\n%s", createBody)
	}
	// The guard must precede the raw call and return the zero value (nil here).
	guardIdx := strings.Index(createBody, "if dataSize < 0 {")
	callIdx := strings.Index(createBody, "capi.CEFBinaryValueCreate")
	if guardIdx < 0 || callIdx < 0 || guardIdx > callIdx {
		t.Errorf("guard must appear before the capi call\n\nGot:\n%s", createBody)
	}
	if !strings.Contains(createBody[guardIdx:callIdx], "return nil") {
		t.Errorf("guard must return the zero value (nil)\n\nGot:\n%s", createBody)
	}
}

// TestEmitPortOutSizeParamUsesUintptr verifies that size_t free-function params
// are mapped to uintptr (matching the capi boundary) in the outbound port
// interface, while pointer params remain unsafe.Pointer.
func TestEmitPortOutSizeParamUsesUintptr(t *testing.T) {
	header := &model.Header{
		Structs: []model.Struct{{
			CName:  "_cef_binary_value_t",
			GoName: "CEFBinaryValueT",
			Kind:   "object",
			Fields: []model.Field{
				{CName: "base", GoName: "Base", CType: "cef_base_ref_counted_t"},
				{
					CName:       "is_valid",
					GoName:      "IsValid",
					IsFunction:  true,
					ReturnCType: "int",
					Params:      []model.Param{{CName: "self", GoName: "self", CType: "struct _cef_binary_value_t*"}},
				},
			},
		}},
		Functions: []model.Function{{
			CName:       "cef_binary_value_create",
			GoName:      "CEFBinaryValueCreate",
			ReturnCType: "struct _cef_binary_value_t*",
			Params: []model.Param{
				{CName: "data", GoName: "data", CType: "const void*", GoType: "unsafe.Pointer"},
				{CName: "data_size", GoName: "dataSize", CType: "size_t", GoType: "uintptr"},
			},
		}},
	}

	registry := NewTypeRegistry([]*model.Header{header})
	portData := BuildPortFileData(header, registry)

	code, err := EmitPortOut(portData)
	if err != nil {
		t.Fatalf("EmitPortOut failed: %v", err)
	}

	if !strings.Contains(code, "BinaryValueCreate(data unsafe.Pointer, dataSize uintptr) unsafe.Pointer") {
		t.Errorf("expected size_t param mapped to uintptr and pointer param kept as unsafe.Pointer\n\nGot:\n%s", code)
	}
}
