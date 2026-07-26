package emitter

import (
	"slices"
	"strings"
)

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
			case "string", "userfreeString", "dataStruct", "outCount", "outSlice", "outObjectSlice":
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

// NeedsMath is retained for template compatibility.
func (d *PublicFileData) NeedsMath() bool {
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

// HasHandlers returns true if there are handler interfaces with methods.
func (d *PublicFileData) HasHandlers() bool {
	for _, iface := range d.Interfaces {
		if iface.Kind == "handler" && len(iface.Methods) > 0 {
			return true
		}
	}
	return false
}

// NeedsRuntime returns true if any generated interface wrapper needs runtime
// finalizers for ref-counted reverse wrappers.
func (d *PublicFileData) NeedsRuntime() bool {
	for _, iface := range d.Interfaces {
		if !iface.IsScoped {
			return true
		}
		for _, method := range iface.Methods {
			for _, param := range method.Params {
				if param.IsRefCountedTransfer {
					return true
				}
			}
		}
	}
	for _, function := range d.FreeFunctions {
		for _, param := range function.Params {
			if param.IsRefCountedTransfer {
				return true
			}
		}
	}
	return false
}

// NeedsSync returns true if generated wrappers need sync.Once for release guards
// or cached typed callback bindings.
func (d *PublicFileData) NeedsSync() bool {
	for _, iface := range d.Interfaces {
		if !iface.IsScoped || slices.ContainsFunc(iface.Methods, methodNeedsTypedObjectCall) {
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
		if slices.ContainsFunc(iface.Methods, methodNeedsTypedObjectCall) {
			return true
		}
	}
	return false
}

// HasSizedDataStructs returns true if any data struct has a Size field
// that needs a New*() constructor with unsafe.Sizeof.
func (d *PublicFileData) HasSizedDataStructs() bool {
	for _, ds := range d.DataStructs {
		if ds.HasSizeField {
			return true
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

func isTypedCallbackPublicType(typ string) bool {
	return isFloatPublicType(typ) || typ == "int64" || typ == "uint64"
}

func methodNeedsTypedObjectCall(m MethodData) bool {
	for _, p := range m.Params {
		if p.MarshalKind == "numeric" && isTypedCallbackPublicType(p.PublicType) {
			return true
		}
	}
	return m.Return.IsNumeric && isTypedCallbackPublicType(m.Return.PublicType)
}

// InterfaceData represents a handler or object interface.
type InterfaceData struct {
	Name      string // "Browser", "LifeSpanHandler"
	Doc       string
	Kind      string // "handler" or "object"
	IsScoped  bool   // true for scoped types (cef_base_scoped_t base)
	RawGoName string // "CEFBrowserT"
	Methods   []MethodData
	// NeedsTake is true when a global factory function returns this ref-counted
	// interface with ownership transferred, requiring a takeX adopter that skips
	// the AddRef performed by wrapX.
	NeedsTake bool
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
	Name                 string // "browser"
	PublicType           string // "Browser"
	CType                string // original C type for marshal decisions
	MarshalKind          string // "interface", "string", "enum", "numeric", "pointer", "userfreeString", "slice", "outSlice", "outObjectSlice", "outCount"
	RawGoType            string // raw Go type for free func params (e.g., "raw.CEFJsonParserOptionsT")
	IsHandler            bool   // true when this param is a handler interface input that can be auto-wrapped
	IsRefCountedTransfer bool   // true when CEF consumes one transferred reference for this interface input

	// For MarshalKind "slice"/"objectSlice": describes how to decode from raw
	// callback args when a count+pointer pair is merged into a single public slice.
	SliceCountArg string // raw arg name for the count (e.g., "arg2")
	SlicePtrArg   string // raw arg name for the pointer (e.g., "arg3")
	SliceElemType string // element type (e.g., "Rect")

	// CountParamName links an out-slice buffer param to its count pointer param
	// (for example, elements []PostDataElement -> elementscount *int).
	CountParamName string
	// CountPartnerName links an out-count param to its slice buffer param.
	CountPartnerName string

	// UnmarshalExtra holds a resolved Go expression used by special marshal kinds
	// (pixelBuffer, objectSlice) that need more context than a single param cast.
	UnmarshalExtra string
}

// NeedsSizeGuard reports whether this param is a signed Go integer passed to the
// C layer as a size_t (uintptr). A negative value would wrap to a huge size_t and
// make libcef read out of bounds, so global wrappers must reject it up front.
func (p ParamData) NeedsSizeGuard() bool {
	if p.RawGoType != "uintptr" {
		return false
	}
	switch p.PublicType {
	case "int", "int8", "int16", "int32", "int64":
		return true
	default:
		return false
	}
}

// PortParamType returns the Go type used for this free-function param in the
// outbound port interface. Size_t and other opaque handle params are uintptr at
// the C boundary (matching the capi signatures); everything else is unsafe.Pointer.
func (p ParamData) PortParamType() string {
	if p.RawGoType == "uintptr" {
		return "uintptr"
	}
	return "unsafe.Pointer"
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
	// UseTake is true when a global function returns this ref-counted interface
	// with ownership transferred; the wrapper must adopt it with takeX (no AddRef)
	// instead of wrapX.
	UseTake bool
}

// DataStructData represents a plain data struct re-export.
type DataStructData struct {
	Name         string
	Doc          string
	RawGoName    string
	Fields       []DataFieldData
	HasSizeField bool // true when the struct has a Size uintptr field (needs New*() constructor)
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
