package emitter

import (
	"testing"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

func testRegistry() *TypeRegistry {
	headers := []*model.Header{
		{
			Structs: []model.Struct{
				{CName: "_cef_browser_t", Kind: "object", InterfaceName: "Browser"},
				{CName: "cef_browser_t", Kind: "object", InterfaceName: "Browser"},
				{CName: "_cef_life_span_handler_t", Kind: "handler", InterfaceName: "LifeSpanHandler"},
				{CName: "cef_life_span_handler_t", Kind: "handler", InterfaceName: "LifeSpanHandler"},
				{CName: "_cef_request_t", Kind: "object", InterfaceName: "Request"},
				{CName: "cef_request_t", Kind: "object", InterfaceName: "Request"},
				{CName: "_cef_popup_features_t", Kind: "data", InterfaceName: "PopupFeatures"},
				{CName: "cef_popup_features_t", Kind: "data", InterfaceName: "PopupFeatures"},
				{CName: "_cef_settings_t", Kind: "data", InterfaceName: "Settings"},
				{CName: "cef_settings_t", Kind: "data", InterfaceName: "Settings"},
				{CName: "_cef_main_args_t", Kind: "data", InterfaceName: "MainArgs"},
				{CName: "cef_main_args_t", Kind: "data", InterfaceName: "MainArgs"},
			},
			Enums: []model.Enum{
				{CName: "cef_process_id_t"},
				{CName: "cef_log_severity_t"},
			},
		},
	}
	return NewTypeRegistry(headers)
}

func TestResolvePublicType(t *testing.T) {
	r := testRegistry()

	tests := []struct {
		ctype string
		want  string
	}{
		// void
		{"void", ""},
		// void pointers
		{"void*", "unsafe.Pointer"},
		{"const void*", "unsafe.Pointer"},
		// string types
		{"const cef_string_t*", "string"},
		{"cef_string_t*", "uintptr"},
		{"cef_string_userfree_t", "string"},
		// string collections (opaque handles)
		{"cef_string_list_t", "StringList"},
		{"cef_string_map_t", "StringMap"},
		{"cef_string_multimap_t", "StringMultimap"},
		// primitive types
		{"int", "int32"},
		{"unsigned int", "uint32"},
		{"size_t", "int"},
		{"double", "float64"},
		{"float", "float32"},
		// stdint types
		{"int8_t", "int8"},
		{"uint8_t", "uint8"},
		{"int16_t", "int16"},
		{"uint16_t", "uint16"},
		{"int32_t", "int32"},
		{"uint32_t", "uint32"},
		{"int64_t", "int64"},
		{"uint64_t", "uint64"},
		// char pointer
		{"const char*", "string"},
		{"char*", "string"},
		// typed numeric pointers
		{"size_t*", "*int"},
		{"int*", "*int32"},
		{"uint32_t*", "*uint32"},
		{"int64_t*", "*int64"},
		{"uint64_t*", "*uint64"},
		{"float*", "*float32"},
		{"double*", "*float64"},
		// struct pointer: object kind -> interface name
		{"struct _cef_browser_t*", "Browser"},
		// struct pointer: handler kind -> raw handler name
		{"struct _cef_life_span_handler_t*", "RawLifeSpanHandler"},
		// struct pointer: data kind -> pointer to data struct
		{"struct _cef_popup_features_t*", "*PopupFeatures"},
		// const struct pointer
		{"const struct _cef_popup_features_t*", "*PopupFeatures"},
		// enum type
		{"cef_process_id_t", "ProcessID"},
		{"cef_log_severity_t", "LogSeverity"},
		// renamed public types
		{"struct _cef_settings_t*", "*RawSettings"},
		{"struct _cef_main_args_t*", "*RawMainArgs"},
		// "TYPE const*" form (same as "const TYPE*")
		{"cef_popup_features_t const*", "*PopupFeatures"},
		{"struct _cef_browser_t const*", "Browser"},
		// bare pointer (no struct keyword)
		{"cef_browser_t*", "Browser"},
		{"cef_popup_features_t*", "*PopupFeatures"},
		// fallback
		{"some_unknown_type", "uintptr"},
	}

	for _, tt := range tests {
		t.Run(tt.ctype, func(t *testing.T) {
			got := r.ResolvePublicType(tt.ctype)
			if got != tt.want {
				t.Errorf("ResolvePublicType(%q) = %q, want %q", tt.ctype, got, tt.want)
			}
		})
	}
}

func TestIsGetterCallback(t *testing.T) {
	r := testRegistry()

	selfParam := model.Param{CName: "self", CType: "struct _cef_life_span_handler_t*"}

	tests := []struct {
		name string
		f    model.Field
		want bool
	}{
		{
			name: "getter returning handler pointer with self only",
			f: model.Field{
				CName:       "get_life_span_handler",
				IsFunction:  true,
				Params:      []model.Param{selfParam},
				ReturnCType: "struct _cef_life_span_handler_t*",
			},
			want: true,
		},
		{
			name: "not a function field",
			f: model.Field{
				CName:      "some_field",
				IsFunction: false,
			},
			want: false,
		},
		{
			name: "function with extra params",
			f: model.Field{
				CName:      "get_something",
				IsFunction: true,
				Params: []model.Param{
					selfParam,
					{CName: "index", CType: "int"},
				},
				ReturnCType: "struct _cef_life_span_handler_t*",
			},
			want: false,
		},
		{
			name: "function returning object (not handler)",
			f: model.Field{
				CName:       "get_browser",
				IsFunction:  true,
				Params:      []model.Param{selfParam},
				ReturnCType: "struct _cef_browser_t*",
			},
			want: false,
		},
		{
			name: "function returning void",
			f: model.Field{
				CName:       "do_something",
				IsFunction:  true,
				Params:      []model.Param{selfParam},
				ReturnCType: "void",
			},
			want: false,
		},
		{
			name: "function returning int",
			f: model.Field{
				CName:       "get_count",
				IsFunction:  true,
				Params:      []model.Param{selfParam},
				ReturnCType: "int",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.IsGetterCallback(tt.f)
			if got != tt.want {
				t.Errorf("IsGetterCallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBoolReturn(t *testing.T) {
	tests := []struct {
		name string
		f    model.Field
		want bool
	}{
		{
			name: "is_ prefix with int return",
			f: model.Field{
				CName:       "is_visible",
				IsFunction:  true,
				ReturnCType: "int",
			},
			want: true,
		},
		{
			name: "has_ prefix with int return",
			f: model.Field{
				CName:       "has_document",
				IsFunction:  true,
				ReturnCType: "int",
			},
			want: true,
		},
		{
			name: "can_ prefix with int return",
			f: model.Field{
				CName:       "can_go_back",
				IsFunction:  true,
				ReturnCType: "int",
			},
			want: true,
		},
		{
			name: "do_ prefix with int return",
			f: model.Field{
				CName:       "do_close",
				IsFunction:  true,
				ReturnCType: "int",
			},
			want: true,
		},
		{
			name: "on_before_ prefix with int return",
			f: model.Field{
				CName:       "on_before_popup",
				IsFunction:  true,
				ReturnCType: "int",
			},
			want: true,
		},
		{
			name: "not a function",
			f: model.Field{
				CName:       "is_visible",
				IsFunction:  false,
				ReturnCType: "int",
			},
			want: false,
		},
		{
			name: "returns void not int",
			f: model.Field{
				CName:       "is_visible",
				IsFunction:  true,
				ReturnCType: "void",
			},
			want: false,
		},
		{
			name: "no bool prefix",
			f: model.Field{
				CName:       "get_count",
				IsFunction:  true,
				ReturnCType: "int",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBoolReturn(tt.f)
			if got != tt.want {
				t.Errorf("IsBoolReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeConst(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"const cef_rect_t*", "const cef_rect_t*"}, // already normalized
		{"cef_rect_t const*", "const cef_rect_t*"}, // needs normalization
		{"struct _cef_rect_t const*", "const struct _cef_rect_t*"},
		{"cef_rect_t*", "cef_rect_t*"},                           // no const, unchanged
		{"const void*", "const void*"},                           // already normalized
		{"void*", "void*"},                                       // no const
		{"int", "int"},                                           // no pointer, no const
		{"cef_rect_t const* const*", "const cef_rect_t const**"}, // double pointer — inner const preserved
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeConst(tt.input)
			if got != tt.want {
				t.Errorf("normalizeConst(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewTypeRegistry(t *testing.T) {
	r := testRegistry()

	if len(r.structs) == 0 {
		t.Error("expected structs in registry")
	}
	if len(r.enums) == 0 {
		t.Error("expected enums in registry")
	}

	// Verify struct lookup works.
	if s, ok := r.structs["_cef_browser_t"]; !ok {
		t.Error("expected _cef_browser_t in registry")
	} else if s.Kind != "object" {
		t.Errorf("expected kind 'object', got %q", s.Kind)
	}

	// Verify enum lookup works.
	if _, ok := r.enums["cef_process_id_t"]; !ok {
		t.Error("expected cef_process_id_t in registry")
	}
}
