package emitter

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// EmitPublic takes a PublicFileData view model and returns formatted Go source
// for the public API layer.
func EmitPublic(data *PublicFileData) (string, error) {
	isFloatPublicType := func(typ string) bool {
		return typ == "float32" || typ == "float64"
	}
	constructorName := func(typeName string) string {
		if rawConstructorTypes[typeName] {
			return "newRaw" + typeName
		}
		return "New" + typeName
	}
	interfaceRawPointerExpr := func(p ParamData) string {
		if p.IsHandler {
			return "extractOrWrapRawPointer(" + p.Name + ", func() any { return " + constructorName(p.PublicType) + "(" + p.Name + ") })"
		}
		return "extractRawPointer(" + p.Name + ")"
	}
	marshalNonFloatAsUintptr := func(p ParamData) string {
		switch p.MarshalKind {
		case "interface":
			if p.PublicType == "unsafe.Pointer" {
				return "uintptr(" + p.Name + ")"
			}
			return "uintptr(" + interfaceRawPointerExpr(p) + ")"
		case "string", "userfreeString":
			return "uintptr(unsafe.Pointer(&" + p.Name + "Str))"
		case "enum":
			return "uintptr(" + p.Name + ")"
		case "pointer":
			return "uintptr(" + p.Name + ")"
		case "dataStruct":
			return "uintptr(unsafe.Pointer(" + p.Name + "))"
		case "outCount":
			return "uintptr(unsafe.Pointer(" + p.CountPartnerName + "CountPtr))"
		case "outSlice", "outObjectSlice":
			return "uintptr(" + p.Name + "Ptr)"
		case "numeric":
			if p.PublicType == "uintptr" {
				return p.Name
			}
			return "uintptr(" + p.Name + ")"
		default:
			return "uintptr(" + p.Name + ")"
		}
	}

	funcMap := template.FuncMap{
		"lower": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		// constructorName returns "newRawX" for types with handwritten public
		// constructors in bridge.go, "NewX" for everything else.
		"constructorName": constructorName,
		"zeroVal": func(typ string) string {
			switch typ {
			case "bool":
				return "false"
			case "string":
				return `""`
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"uintptr", "float32", "float64":
				return "0"
			default:
				if strings.HasPrefix(typ, "*") || strings.HasPrefix(typ, "[]") ||
					strings.HasPrefix(typ, "map[") || typ == "unsafe.Pointer" {
					return "nil"
				}
				// Interface types default to nil.
				return "nil"
			}
		},
		"isEnumReturn": func(ret ReturnData) bool {
			return ret.IsEnum
		},
		"isFloatPublicType": isFloatPublicType,
		"needsTypedObjectCall": func(m MethodData) bool {
			return methodNeedsTypedObjectCall(m)
		},
		"typedObjectFuncSig": func(rawGoName string, m MethodData) string {
			parts := []string{"*capi." + rawGoName}
			for _, p := range m.Params {
				switch {
				case p.MarshalKind == "slice", p.MarshalKind == "objectSlice":
					parts = append(parts, "uintptr", "uintptr")
				case p.MarshalKind == "numeric" && isFloatPublicType(p.PublicType):
					parts = append(parts, p.PublicType)
				default:
					parts = append(parts, "uintptr")
				}
			}
			sig := "func(" + strings.Join(parts, ", ") + ")"
			if !m.Return.IsVoid {
				if m.Return.IsNumeric && isFloatPublicType(m.Return.PublicType) {
					sig += " " + m.Return.PublicType
				} else {
					sig += " uintptr"
				}
			}
			return sig
		},
		"typedObjectCallArgs": func(m MethodData) string {
			args := []string{"obj.rawPtr"}
			for _, p := range m.Params {
				if p.MarshalKind == "slice" || p.MarshalKind == "objectSlice" {
					args = append(args, "uintptr(len("+p.Name+"))", "uintptr("+p.Name+"Ptr)")
					continue
				}
				if p.MarshalKind == "numeric" && isFloatPublicType(p.PublicType) {
					args = append(args, p.Name)
					continue
				}
				args = append(args, marshalNonFloatAsUintptr(p))
			}
			return strings.Join(args, ", ")
		},
		// cbParamName generates a raw parameter name for callbacks (e.g., "arg0", "arg1").
		"cbParamName": func(p ParamData, idx int) string {
			return fmt.Sprintf("arg%d", idx)
		},
		// rawArgForParam finds the raw callback arg name for a public param by matching names.
		"rawArgForParam": func(p ParamData, rawParams []ParamData) string {
			for i, rp := range rawParams {
				if rp.Name == p.Name {
					return fmt.Sprintf("arg%d", i)
				}
			}
			panic(fmt.Sprintf("rawArgForParam: no raw param match for %q", p.Name))
		},
		// sliceCountArg finds the raw arg for the count param of a merged slice.
		"sliceCountArg": func(p ParamData, rawParams []ParamData) string {
			// The count param precedes the pointer param and has name = ptrName + "count"
			for i, rp := range rawParams {
				if strings.ToLower(rp.Name) == strings.ToLower(p.Name)+"count" {
					return fmt.Sprintf("arg%d", i)
				}
			}
			return "0"
		},
		// slicePtrArg finds the raw arg for the pointer param of a merged slice.
		"slicePtrArg": func(p ParamData, rawParams []ParamData) string {
			for i, rp := range rawParams {
				if rp.Name == p.Name {
					return fmt.Sprintf("arg%d", i)
				}
			}
			return "0"
		},
		// marshalPreamble returns pre-call statements needed before the Call (e.g., string setup).
		// Returns empty string if no preamble is needed.
		"marshalPreamble": func(p ParamData) string {
			switch p.MarshalKind {
			case "string", "userfreeString":
				return p.Name + "Str := cefString(" + p.Name + ")\n\tdefer freeCefString(&" + p.Name + "Str)"
			case "slice":
				return "var " + p.Name + "Ptr unsafe.Pointer\n\tif len(" + p.Name + ") > 0 {\n\t\t" + p.Name + "Ptr = unsafe.Pointer(&" + p.Name + "[0])\n\t}"
			case "objectSlice":
				return "var " + p.Name + "Raw []uintptr\n\tvar " + p.Name + "Ptr unsafe.Pointer\n\tif len(" + p.Name + ") > 0 {\n\t\t" + p.Name + "Raw = make([]uintptr, len(" + p.Name + "))\n\t\tfor i, elem := range " + p.Name + " {\n\t\t\t" + p.Name + "Raw[i] = uintptr(extractRawPointer(elem))\n\t\t}\n\t\t" + p.Name + "Ptr = unsafe.Pointer(&" + p.Name + "Raw[0])\n\t}"
			case "outSlice":
				return "var " + p.Name + "Ptr unsafe.Pointer\n\t" + p.Name + "CountPtr := " + p.CountParamName + "\n\tif " + p.Name + "CountPtr == nil {\n\t\t" + p.Name + "CountScratch := len(" + p.Name + ")\n\t\t" + p.Name + "CountPtr = &" + p.Name + "CountScratch\n\t} else {\n\t\t*" + p.Name + "CountPtr = len(" + p.Name + ")\n\t}\n\tif len(" + p.Name + ") > 0 {\n\t\t" + p.Name + "Ptr = unsafe.Pointer(&" + p.Name + "[0])\n\t}"
			case "outObjectSlice":
				return "var " + p.Name + "Raw []uintptr\n\tvar " + p.Name + "Ptr unsafe.Pointer\n\t" + p.Name + "CountPtr := " + p.CountParamName + "\n\tif " + p.Name + "CountPtr == nil {\n\t\t" + p.Name + "CountScratch := len(" + p.Name + ")\n\t\t" + p.Name + "CountPtr = &" + p.Name + "CountScratch\n\t} else {\n\t\t*" + p.Name + "CountPtr = len(" + p.Name + ")\n\t}\n\tif len(" + p.Name + ") > 0 {\n\t\t" + p.Name + "Raw = make([]uintptr, len(" + p.Name + "))\n\t\t" + p.Name + "Ptr = unsafe.Pointer(&" + p.Name + "Raw[0])\n\t}"
			default:
				return ""
			}
		},
		"marshalPostamble": func(p ParamData) string {
			switch p.MarshalKind {
			case "outObjectSlice":
				return "if len(" + p.Name + ") > 0 {\n\t\tn := len(" + p.Name + ")\n\t\tif *" + p.Name + "CountPtr < n {\n\t\t\tn = *" + p.Name + "CountPtr\n\t\t}\n\t\tfor i := 0; i < n; i++ {\n\t\t\t" + p.Name + "[i] = wrap" + p.SliceElemType + "(unsafe.Pointer(" + p.Name + "Raw[i]))\n\t\t}\n\t}"
			default:
				return ""
			}
		},
		"hasMarshalPostamble": func(params []ParamData) bool {
			for _, p := range params {
				if p.MarshalKind == "outObjectSlice" {
					return true
				}
			}
			return false
		},
		// methodPreamble returns method-specific pre-call statements for public object wrappers.
		"methodPreamble": func(ifaceName string, m MethodData) string {
			switch ifaceName + "." + m.Name {
			case "V8Context.Eval":
				return "\t// CEF's eval requires valid output pointers — it unconditionally writes\n" +
					"\t// *retval and *exception. Provide scratch storage when the caller passes nil\n" +
					"\t// to avoid a NULL-pointer dereference crash in the renderer subprocess.\n" +
					"\tif retval == nil {\n" +
					"\t\tvar scratch uintptr\n" +
					"\t\tretval = unsafe.Pointer(&scratch)\n" +
					"\t}\n" +
					"\tif exception == nil {\n" +
					"\t\tvar scratch uintptr\n" +
					"\t\texception = unsafe.Pointer(&scratch)\n" +
					"\t}\n"
			default:
				return ""
			}
		},
		// marshalCallArgs generates the full comma-separated argument list for a raw Call,
		// expanding slice params into count + pointer pairs.
		"marshalCallArgs": func(params []ParamData) string {
			marshalOne := func(p ParamData) string {
				switch p.MarshalKind {
				case "interface":
					if p.PublicType == "unsafe.Pointer" {
						return "uintptr(" + p.Name + ")"
					}
					return "uintptr(" + interfaceRawPointerExpr(p) + ")"
				case "string", "userfreeString":
					return "uintptr(unsafe.Pointer(&" + p.Name + "Str))"
				case "enum":
					return "uintptr(" + p.Name + ")"
				case "pointer":
					return "uintptr(" + p.Name + ")"
				case "dataStruct":
					return "uintptr(unsafe.Pointer(" + p.Name + "))"
				case "outCount":
					return "uintptr(unsafe.Pointer(" + p.CountPartnerName + "CountPtr))"
				case "outSlice", "outObjectSlice":
					return "uintptr(" + p.Name + "Ptr)"
				case "numeric":
					switch p.PublicType {
					case "float64":
						return "uintptr(math.Float64bits(" + p.Name + "))"
					case "float32":
						return "uintptr(math.Float32bits(" + p.Name + "))"
					case "uintptr":
						return p.Name
					default:
						return "uintptr(" + p.Name + ")"
					}
				default:
					return "uintptr(" + p.Name + ")"
				}
			}

			var args []string
			for _, p := range params {
				if p.MarshalKind == "slice" || p.MarshalKind == "objectSlice" {
					args = append(args, "uintptr(len("+p.Name+"))")
					args = append(args, "uintptr("+p.Name+"Ptr)")
				} else {
					args = append(args, marshalOne(p))
				}
			}
			return strings.Join(args, ", ")
		},
		// marshalParamForRawFunc generates the Go expression for free function calls.
		// Raw free functions use typed params (unsafe.Pointer, int32, etc.), not uintptr.
		"marshalParamForRawFunc": func(p ParamData) string {
			switch p.MarshalKind {
			case "interface":
				return interfaceRawPointerExpr(p)
			case "string", "userfreeString":
				return "unsafe.Pointer(&" + p.Name + "Str)"
			case "dataStruct":
				return "unsafe.Pointer(" + p.Name + ")"
			case "outCount":
				return "unsafe.Pointer(" + p.CountPartnerName + "CountPtr)"
			case "outSlice", "outObjectSlice":
				return p.Name + "Ptr"
			case "enum":
				// Cast public enum type to raw enum type if they differ.
				if p.RawGoType != "" && p.RawGoType != "unsafe.Pointer" {
					return "capi." + p.RawGoType + "(" + p.Name + ")"
				}
				return p.Name
			case "numeric":
				// Cast to raw type if it differs from public type.
				if p.RawGoType != "" && p.RawGoType != p.PublicType {
					// CEF types live in raw package.
					if strings.HasPrefix(p.RawGoType, "CEF") {
						return "capi." + p.RawGoType + "(" + p.Name + ")"
					}
					return p.RawGoType + "(" + p.Name + ")"
				}
				return p.Name
			default:
				// pointer — pass directly
				return p.Name
			}
		},
		// marshalParam generates the Go expression to convert a public Go param to uintptr.
		"marshalParam": func(p ParamData) string {
			switch p.MarshalKind {
			case "interface":
				if p.PublicType == "unsafe.Pointer" {
					return "uintptr(" + p.Name + ")"
				}
				return "uintptr(" + interfaceRawPointerExpr(p) + ")"
			case "string", "userfreeString":
				return "uintptr(unsafe.Pointer(&" + p.Name + "Str))"
			case "enum":
				return "uintptr(" + p.Name + ")"
			case "pointer":
				return "uintptr(" + p.Name + ")"
			case "dataStruct":
				return "uintptr(unsafe.Pointer(" + p.Name + "))"
			case "outCount":
				return "uintptr(unsafe.Pointer(" + p.CountPartnerName + "CountPtr))"
			case "outSlice", "outObjectSlice":
				return "uintptr(" + p.Name + "Ptr)"
			case "numeric":
				switch p.PublicType {
				case "float64":
					return "uintptr(math.Float64bits(" + p.Name + "))"
				case "float32":
					return "uintptr(math.Float32bits(" + p.Name + "))"
				case "uintptr":
					return p.Name
				default:
					return "uintptr(" + p.Name + ")"
				}
			default:
				return "uintptr(" + p.Name + ")"
			}
		},
		// unmarshalObjectSlice generates multi-line code to unmarshal an objectSlice param.
		"unmarshalObjectSlice": func(p ParamData, rawParams []ParamData) string {
			countArg := "0"
			ptrArg := "0"
			ptrName := strings.ToLower(p.Name)
			for i, rp := range rawParams {
				rpName := strings.ToLower(rp.Name)
				if rpName == ptrName+"count" {
					countArg = fmt.Sprintf("arg%d", i)
				}
				if rpName == ptrName {
					ptrArg = fmt.Sprintf("arg%d", i)
				}
			}
			elemType := p.SliceElemType
			var b strings.Builder
			fmt.Fprintf(&b, "var %s []%s\n", p.Name, elemType)
			fmt.Fprintf(&b, "\t\tif %s != 0 && %s > 0 {\n", ptrArg, countArg)
			fmt.Fprintf(&b, "\t\t\t%sPtrs := unsafe.Slice((*uintptr)(unsafe.Pointer(%s)), int(%s))\n", p.Name, ptrArg, countArg)
			fmt.Fprintf(&b, "\t\t\t%s = make([]%s, int(%s))\n", p.Name, elemType, countArg)
			fmt.Fprintf(&b, "\t\t\tfor i, ptr := range %sPtrs {\n", p.Name)
			fmt.Fprintf(&b, "\t\t\t\t%s[i] = wrap%s(unsafe.Pointer(ptr))\n", p.Name, elemType)
			fmt.Fprintf(&b, "\t\t\t}\n")
			fmt.Fprintf(&b, "\t\t}")
			return b.String()
		},
		// unmarshalParam generates the Go expression to convert a uintptr to the public type.
		"unmarshalParam": func(p ParamData, rawName string) string {
			switch p.MarshalKind {
			case "interface":
				if p.PublicType == "unsafe.Pointer" {
					return "unsafe.Pointer(" + rawName + ")"
				}
				return "wrap" + p.PublicType + "(unsafe.Pointer(" + rawName + "))"
			case "string":
				return "goString(unsafe.Pointer(" + rawName + "))"
			case "userfreeString":
				return "goStringUserfree(unsafe.Pointer(" + rawName + "))"
			case "enum":
				return p.PublicType + "(" + rawName + ")"
			case "pointer":
				return "unsafe.Pointer(" + rawName + ")"
			case "pixelBuffer":
				return "unsafe.Slice((*byte)(unsafe.Pointer(" + rawName + ")), " + p.UnmarshalExtra + ")"
			case "dataStruct":
				// PublicType is "*WindowInfo" etc, we need "(*WindowInfo)(unsafe.Pointer(...))"
				return "(" + p.PublicType + ")(unsafe.Pointer(" + rawName + "))"
			case "numeric":
				if isFloatPublicType(p.PublicType) {
					return rawName
				}
				return p.PublicType + "(" + rawName + ")"
			default:
				return p.PublicType + "(" + rawName + ")"
			}
		},
	}

	tmpl, err := template.New("public").Funcs(funcMap).ParseFS(templateFS,
		"templates/public_file.tmpl",
		"templates/interface.tmpl",
		"templates/object_wrapper.tmpl",
		"templates/handler_constructor.tmpl",
		"templates/data_struct.tmpl",
		"templates/enums.tmpl",
		"templates/free_func.tmpl",
	)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "public_file.tmpl", data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}

// EmitRaw takes a parsed Header and returns formatted Go source for the raw layer.
// The raw layer lives in package "raw" and mirrors C struct layouts.
func EmitRaw(header *model.Header) (string, error) {
	tmpl, err := template.New("raw").ParseFS(templateFS, "templates/raw_*.tmpl")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "raw_file.tmpl", header); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}

// EmitPortOut generates outbound port interfaces for a header.
// This includes per-object interfaces AND a Functions interface for free functions.
func EmitPortOut(data *PublicFileData) (string, error) {
	if data == nil {
		return "", nil
	}
	hasObjects := false
	for _, iface := range data.Interfaces {
		if iface.Kind == "object" {
			hasObjects = true
			break
		}
	}
	hasFreeFuncs := len(data.FreeFunctions) > 0
	if !hasObjects && !hasFreeFuncs {
		return "", nil
	}
	tmpl, err := template.New("port_out").ParseFS(templateFS, "templates/port_out.tmpl")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "port_out.tmpl", data); err != nil {
		return "", fmt.Errorf("execute port_out template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format port_out source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}

// EmitPortIn generates inbound port interfaces for a header.
func EmitPortIn(data *PublicFileData) (string, error) {
	if data == nil {
		return "", nil
	}
	hasContent := len(data.Interfaces) > 0 || len(data.Enums) > 0 || len(data.DataStructs) > 0
	if !hasContent {
		return "", nil
	}
	tmpl, err := template.New("port_in").ParseFS(templateFS, "templates/port_in.tmpl")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "port_in.tmpl", data); err != nil {
		return "", fmt.Errorf("execute port_in template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format port_in source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}
