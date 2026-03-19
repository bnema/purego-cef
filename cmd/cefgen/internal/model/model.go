package model

import "strings"

type File struct {
	PackageName string
	Headers     []string
	Structs     []Struct
	Functions   []Function
	Enums       []Enum
}

type Struct struct {
	CName         string
	GoName        string
	Kind          string // "handler", "object", or "data"
	InterfaceName string // public Go name (e.g., "Browser")
	Doc           string // extracted doc comment
	Fields        []Field
	Comments      []string
}

type Field struct {
	CName                 string
	GoName                string
	CType                 string
	GoType                string
	Doc                   string // extracted doc comment
	GoInterfaceType       string // resolved public type for params
	ReturnGoInterfaceType string // resolved public type for returns
	IsFunction            bool
	IsPointer             bool
	Params                []Param
	ReturnCType           string
	ReturnGoType          string
	Comments              []string
}

type Function struct {
	CName        string
	GoName       string
	Doc          string // extracted doc comment
	Params       []Param
	ReturnCType  string
	ReturnGoType string
	Comments     []string
}

type Param struct {
	CName           string
	GoName          string
	CType           string
	GoType          string
	GoInterfaceType string // resolved public type for params
	IsConst         bool
	Pointer         int
}

type Enum struct {
	CName    string
	GoName   string
	Values   []EnumValue
	Comments []string
}

type EnumValue struct {
	CName    string
	GoName   string
	Value    string
	Comments []string
}

type Header struct {
	Path         string
	Package      string
	RegisterName string // e.g. "RegisterApp" for cef_app_capi.h
	Structs      []Struct
	Functions    []Function
	Enums        []Enum
	Typedefs     []Typedef
}

type Typedef struct {
	CName  string
	GoName string
	CType  string
	GoType string
}

// NeedsStructs returns true if the header has any struct definitions.
func (h *Header) NeedsStructs() bool {
	return len(h.Structs) > 0
}

// NeedsUnsafe returns true if any field or function parameter uses unsafe.Pointer.
func (h *Header) NeedsUnsafe() bool {
	for _, s := range h.Structs {
		for _, f := range s.Fields {
			if f.GoType == "unsafe.Pointer" {
				return true
			}
			if f.ReturnGoType == "unsafe.Pointer" {
				return true
			}
			for _, p := range f.Params {
				if p.GoType == "unsafe.Pointer" {
					return true
				}
			}
		}
	}
	for _, fn := range h.Functions {
		if fn.ReturnGoType == "unsafe.Pointer" {
			return true
		}
		for _, p := range fn.Params {
			if p.GoType == "unsafe.Pointer" {
				return true
			}
		}
	}
	return false
}

// NeedsPurego returns true if the header has any free functions to register
// or any struct with function pointer fields (which generate Call methods using purego.SyscallN).
func (h *Header) NeedsPurego() bool {
	if len(h.Functions) > 0 {
		return true
	}
	for _, s := range h.Structs {
		for _, f := range s.Fields {
			if f.IsFunction {
				return true
			}
		}
	}
	return false
}

// acronyms maps lowercase words to their uppercase Go equivalents.
var acronyms = map[string]string{
	"id":  "ID",
	"url": "URL",
	"ui":  "UI",
}

// PublicName converts a C type name to a clean public Go name.
// It strips leading underscores, the "cef_" prefix, and the "_t" suffix,
// then PascalCases the remaining underscore-separated segments.
//
// Examples:
//
//	cef_browser_t        → Browser
//	cef_life_span_handler_t → LifeSpanHandler
//	_cef_browser_t       → Browser
//	cef_request_context_t → RequestContext
func PublicName(cName string) string {
	s := strings.TrimLeft(cName, "_")
	s = strings.TrimPrefix(s, "cef_")
	s = strings.TrimSuffix(s, "_t")

	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if upper, ok := acronyms[strings.ToLower(p)]; ok {
			b.WriteString(upper)
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return b.String()
}
