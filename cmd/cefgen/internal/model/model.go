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
	Path      string
	Package   string
	Structs   []Struct
	Functions []Function
	Enums     []Enum
	Typedefs  []Typedef
}

type Typedef struct {
	CName  string
	GoName string
	CType  string
	GoType string
}
