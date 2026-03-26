package emitter

import "strings"

// PublicFileData holds all data needed to render the public API file for one header.
type PublicFileData struct {
	PackageName   string
	HeaderStem    string // e.g., "App", "Browser" — used for naming generated interfaces
	Interfaces    []InterfaceData
	DataStructs   []DataStructData
	Enums         []EnumData
	FreeFunctions []FreeFuncData
}

// NeedsUnsafe returns true if any type reference uses unsafe.Pointer.
func (d *PublicFileData) NeedsUnsafe() bool {
	for _, iface := range d.Interfaces {
		for _, m := range iface.Methods {
			if usesUnsafe(m.Return.PublicType) {
				return true
			}
			for _, p := range m.Params {
				if usesUnsafe(p.PublicType) {
					return true
				}
			}
		}
	}
	for _, ff := range d.FreeFunctions {
		if usesUnsafe(ff.Return.PublicType) {
			return true
		}
		for _, p := range ff.Params {
			if usesUnsafe(p.PublicType) {
				return true
			}
			// Free func marshaling generates unsafe.Pointer for these kinds.
			switch p.MarshalKind {
			case "string", "userfreeString", "dataStruct":
				return true
			case "numeric":
				if p.PublicType == "uintptr" {
					return true
				}
			}
		}
	}
	return false
}

func usesUnsafe(typ string) bool {
	return strings.Contains(typ, "unsafe.Pointer")
}

// NeedsUnsafeInSignatures returns true if any interface method signature
// (params or return) directly uses unsafe.Pointer. Unlike NeedsUnsafe(),
// this does NOT check free function marshaling — only literal type references.
func (d *PublicFileData) NeedsUnsafeInSignatures() bool {
	for _, iface := range d.Interfaces {
		for _, m := range iface.Methods {
			if usesUnsafe(m.Return.PublicType) {
				return true
			}
			for _, p := range m.Params {
				if usesUnsafe(p.PublicType) {
					return true
				}
			}
		}
	}
	return false
}

// NeedsMath returns true if any generated callback unmarshalling requires
// math.Float64frombits/Float32frombits.
func (d *PublicFileData) NeedsMath() bool {
	for _, iface := range d.Interfaces {
		if iface.Kind != "handler" {
			continue
		}
		for _, m := range iface.Methods {
			for _, p := range m.Params {
				if p.MarshalKind == "numeric" && (p.PublicType == "float64" || p.PublicType == "float32") {
					return true
				}
			}
		}
	}
	return false
}

// HasObjects returns true if there are non-scoped object interfaces.
func (d *PublicFileData) HasObjects() bool {
	for _, iface := range d.Interfaces {
		if iface.Kind == "object" && !iface.IsScoped {
			return true
		}
	}
	return false
}

// HasHandlers returns true if there are handler interfaces with methods
// (which need purego.NewCallback).
func (d *PublicFileData) HasHandlers() bool {
	for _, iface := range d.Interfaces {
		if iface.Kind == "handler" && len(iface.Methods) > 0 {
			return true
		}
	}
	return false
}

// NeedsPuregoObjectCalls returns true if any object method needs a typed
// purego.RegisterFunc binding to preserve float ABI semantics.
func (d *PublicFileData) NeedsPuregoObjectCalls() bool {
	for _, iface := range d.Interfaces {
		if iface.Kind != "object" {
			continue
		}
		for _, m := range iface.Methods {
			if methodNeedsTypedObjectCall(m) {
				return true
			}
		}
	}
	return false
}

// NeedsRaw returns true if the file references the raw package.
func (d *PublicFileData) NeedsRaw() bool {
	return len(d.Interfaces) > 0 || len(d.DataStructs) > 0 || len(d.FreeFunctions) > 0
}

func isFloatPublicType(typ string) bool {
	return typ == "float32" || typ == "float64"
}

func methodNeedsTypedObjectCall(m MethodData) bool {
	for _, p := range m.Params {
		if p.MarshalKind == "numeric" && isFloatPublicType(p.PublicType) {
			return true
		}
	}
	return m.Return.IsNumeric && isFloatPublicType(m.Return.PublicType)
}

// InterfaceData represents a handler or object interface.
type InterfaceData struct {
	Name      string // "Browser", "LifeSpanHandler"
	Doc       string
	Kind      string // "handler" or "object"
	IsScoped  bool   // true for scoped types (cef_base_scoped_t base)
	RawGoName string // "CEFBrowserT"
	Methods   []MethodData
}

// MethodData represents a single method on an interface.
type MethodData struct {
	Name            string // "GetHost"
	Doc             string
	Params          []ParamData // Public params (may have count+pointer merged into slices)
	RawParams       []ParamData // Original unmerged params (for handler callback signatures)
	Return          ReturnData
	IsGetter        bool   // only self param, returns handler pointer
	GetterInterface string // for getters: "LifeSpanHandler"
	RawFieldName    string // matches raw struct field name
}

// ParamData represents a method parameter.
type ParamData struct {
	Name        string // "browser"
	PublicType  string // "Browser"
	CType       string // original C type for marshal decisions
	MarshalKind string // "interface", "string", "enum", "numeric", "pointer", "userfreeString", "slice"
	RawGoType   string // raw Go type for free func params (e.g., "raw.CEFJsonParserOptionsT")

	// For MarshalKind "slice": describes how to decode from raw callback args.
	SliceCountArg string // raw arg name for the count (e.g., "arg2")
	SlicePtrArg   string // raw arg name for the pointer (e.g., "arg3")
	SliceElemType string // element type (e.g., "Rect")

	// UnmarshalExtra holds a resolved Go expression used by special marshal kinds
	// (pixelBuffer, objectSlice) that need more context than a single param cast.
	UnmarshalExtra string
}

// ReturnData represents a method return type.
type ReturnData struct {
	PublicType   string
	CType        string
	IsBool       bool
	IsVoid       bool
	IsEnum       bool
	IsString     bool // cef_string_userfree_t return
	IsInterface  bool // returns a handler/object interface type
	IsNumeric    bool // int32, int64, uint32, float64, int, etc.
	IsPointer    bool // unsafe.Pointer
	IsHandler    bool // returns a handler interface (needs NewXxx wrapping in callbacks)
	IsDataStruct bool // returns a pointer to a data struct (*Type)
}

// DataStructData represents a plain data struct re-export.
type DataStructData struct {
	Name      string
	Doc       string
	RawGoName string
	Fields    []DataFieldData
}

// DataFieldData represents a field in a data struct.
type DataFieldData struct {
	Name       string
	PublicType string
	Doc        string
}

// EnumData represents an enum type.
type EnumData struct {
	Name      string // "State" (from cef_state_t)
	Doc       string
	RawGoName string // "CEFStateT"
	Unsigned  bool   // true if any value exceeds int32 range
	Values    []EnumValueData
}

// EnumValueData represents a single enum value.
type EnumValueData struct {
	Name  string
	Value string
}

// FreeFuncData represents a free function wrapper.
type FreeFuncData struct {
	Name      string // "Initialize"
	Doc       string
	RawGoName string // "CEFInitialize"
	Params    []ParamData
	Return    ReturnData
}
