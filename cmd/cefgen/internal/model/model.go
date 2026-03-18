package model

type File struct {
	PackageName string
	Headers     []string
	Structs     []Struct
	Functions   []Function
	Enums       []Enum
}

type Struct struct {
	CName    string
	GoName   string
	Fields   []Field
	Comments []string
}

type Field struct {
	CName        string
	GoName       string
	CType        string
	GoType       string
	IsFunction   bool
	IsPointer    bool
	Params       []Param
	ReturnCType  string
	ReturnGoType string
	Comments     []string
}

type Function struct {
	CName        string
	GoName       string
	Params       []Param
	ReturnCType  string
	ReturnGoType string
	Comments     []string
}

type Param struct {
	CName   string
	GoName  string
	CType   string
	GoType  string
	IsConst bool
	Pointer int
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
