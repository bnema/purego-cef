package emitter

// PublicFileData holds all data needed to render the public API file for one header.
type PublicFileData struct {
	PackageName   string
	Interfaces    []InterfaceData
	DataStructs   []DataStructData
	Enums         []EnumData
	FreeFunctions []FreeFuncData
}

// InterfaceData represents a handler or object interface.
type InterfaceData struct {
	Name      string // "Browser", "LifeSpanHandler"
	Doc       string
	Kind      string // "handler" or "object"
	RawGoName string // "CEFBrowserT"
	Methods   []MethodData
}

// MethodData represents a single method on an interface.
type MethodData struct {
	Name            string // "GetHost"
	Doc             string
	Params          []ParamData
	Return          ReturnData
	IsGetter        bool   // only self param, returns handler pointer
	GetterInterface string // for getters: "LifeSpanHandler"
	RawFieldName    string // matches raw struct field name
}

// ParamData represents a method parameter.
type ParamData struct {
	Name       string // "browser"
	PublicType string // "Browser"
	CType      string // original C type for marshal decisions
}

// ReturnData represents a method return type.
type ReturnData struct {
	PublicType string
	CType      string
	IsBool     bool
	IsVoid     bool
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
