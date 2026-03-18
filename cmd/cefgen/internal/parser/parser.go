package parser

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

var (
	structRE   = regexp.MustCompile(`(?s)typedef struct (_cef_[a-z0-9_]+_t)\s*\{(.*?)\}\s*(cef_[a-z0-9_]+_t);`)
	funcRE     = regexp.MustCompile(`CEF_EXPORT\s+(.+?)\s+(cef_[a-z0-9_]+)\((.*?)\);`)
	enumRE     = regexp.MustCompile(`(?s)typedef enum\s*\{(.*?)\}\s*(cef_[a-z0-9_]+_t);`)
	callbackRE = regexp.MustCompile(`(.+?)\(\s*CEF_CALLBACK\*\s*([a-z0-9_]+)\)\((.*?)\)`)
)

// ParseFile reads a CEF capi header file and returns a parsed Header.
func ParseFile(path string) (*model.Header, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, data)
}

// Parse parses the given CEF capi header content.
func Parse(path string, data []byte) (*model.Header, error) {
	clean := stripComments(data)
	// Join continuation lines so multi-line declarations become single lines.
	clean = joinLines(clean)

	out := &model.Header{Path: path, Package: "capi"}
	for _, match := range structRE.FindAllSubmatch(clean, -1) {
		st, err := parseStruct(string(match[1]), string(match[2]), string(match[3]))
		if err != nil {
			return nil, fmt.Errorf("%s: struct %s: %w", path, string(match[3]), err)
		}
		out.Structs = append(out.Structs, st)
	}
	for _, match := range funcRE.FindAllSubmatch(clean, -1) {
		out.Functions = append(out.Functions, parseFunction(string(match[2]), string(match[1]), string(match[3])))
	}
	for _, match := range enumRE.FindAllSubmatch(clean, -1) {
		out.Enums = append(out.Enums, parseEnum(string(match[2]), string(match[1])))
	}
	return out, nil
}

// stripComments removes single-line comments, preprocessor directives, and blank lines.
func stripComments(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	out := make([][]byte, 0, len(lines))
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if bytes.HasPrefix(trimmed, []byte("//")) {
			continue
		}
		if bytes.HasPrefix(trimmed, []byte("#")) {
			continue
		}
		// Strip inline comments
		if idx := bytes.Index(line, []byte("//")); idx >= 0 {
			line = bytes.TrimRightFunc(line[:idx], unicode.IsSpace)
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}
		}
		out = append(out, line)
	}
	return bytes.Join(out, []byte("\n"))
}

// joinLines collapses multi-line declarations into single lines.
// A line that doesn't end with ; or { or } is considered a continuation.
func joinLines(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte
	var buf []byte
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			if buf != nil {
				result = append(result, buf)
				buf = nil
			}
			continue
		}
		if buf == nil {
			buf = bytes.TrimSpace(line)
		} else {
			buf = append(buf, ' ')
			buf = append(buf, bytes.TrimSpace(line)...)
		}
		// Check if the line ends with a statement-ending character
		last := trimmed[len(trimmed)-1]
		if last == ';' || last == '{' || last == '}' {
			result = append(result, buf)
			buf = nil
		}
	}
	if buf != nil {
		result = append(result, buf)
	}
	return bytes.Join(result, []byte("\n"))
}

func parseStruct(tag, body, name string) (model.Struct, error) {
	_ = tag
	result := model.Struct{CName: name, GoName: goName(name)}
	for _, raw := range strings.Split(body, ";") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if match := callbackRE.FindStringSubmatch(line); match != nil {
			result.Fields = append(result.Fields, model.Field{
				CName:       match[2],
				GoName:      goName(match[2]),
				CType:       strings.TrimSpace(match[1]),
				GoType:      "uintptr",
				IsFunction:  true,
				ReturnCType: strings.TrimSpace(match[1]),
				Params:      parseParams(match[3]),
			})
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		fieldName := parts[len(parts)-1]
		fieldType := strings.Join(parts[:len(parts)-1], " ")
		result.Fields = append(result.Fields, model.Field{
			CName:  fieldName,
			GoName: goName(fieldName),
			CType:  fieldType,
			GoType: mapType(fieldType),
		})
	}
	return result, nil
}

func parseParams(raw string) []model.Param {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "void" {
		return nil
	}
	parts := splitParams(raw)
	params := make([]model.Param, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "void" {
			continue
		}
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}
		// The last token is the parameter name (may start with * for pointer)
		name := fields[len(fields)-1]
		// Strip leading * from the name (they belong to the type)
		stars := ""
		for strings.HasPrefix(name, "*") {
			stars += "*"
			name = name[1:]
		}
		ctype := strings.TrimSpace(strings.Join(fields[:len(fields)-1], " ")) + stars
		params = append(params, model.Param{
			CName:   name,
			GoName:  goName(name),
			CType:   ctype,
			GoType:  mapType(ctype),
			IsConst: strings.Contains(ctype, "const "),
			Pointer: strings.Count(ctype, "*"),
		})
	}
	return params
}

// splitParams splits a parameter list by commas, respecting nested parentheses.
func splitParams(raw string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, ch := range raw {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, raw[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, raw[start:])
	return parts
}

func parseFunction(name, ret, params string) model.Function {
	return model.Function{
		CName:        name,
		GoName:       goName(name),
		ReturnCType:  strings.TrimSpace(ret),
		ReturnGoType: mapType(ret),
		Params:       parseParams(params),
	}
}

func parseEnum(name, body string) model.Enum {
	result := model.Enum{CName: name, GoName: goName(name)}
	for _, line := range strings.Split(body, ",") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		cname := strings.TrimSpace(parts[0])
		if cname == "" {
			continue
		}
		val := ""
		if len(parts) == 2 {
			val = strings.TrimSpace(parts[1])
		}
		result.Values = append(result.Values, model.EnumValue{
			CName:  cname,
			GoName: goName(cname),
			Value:  val,
		})
	}
	return result
}

// goName converts a C identifier like cef_client_t or CEF_CALLBACK to a Go name.
func goName(cname string) string {
	// Strip leading underscore (e.g., _cef_client_t -> cef_client_t)
	cname = strings.TrimPrefix(cname, "_")

	// Split on underscores
	parts := strings.Split(cname, "_")
	var sb strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		part = strings.ToLower(part)
		switch part {
		case "cef":
			sb.WriteString("CEF")
		case "t":
			sb.WriteString("T")
		case "id":
			sb.WriteString("ID")
		case "url":
			sb.WriteString("URL")
		case "ui":
			sb.WriteString("UI")
		default:
			sb.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return sb.String()
}

// mapType maps a C type string to a Go type string.
func mapType(ctype string) string {
	ctype = strings.TrimSpace(ctype)

	// Normalize multiple spaces
	ctype = strings.Join(strings.Fields(ctype), " ")

	isPtr := strings.HasSuffix(ctype, "*") || strings.Contains(ctype, "* ")

	// void* -> unsafe.Pointer
	if ctype == "void*" || ctype == "void *" || ctype == "const void*" || ctype == "const void *" {
		return "unsafe.Pointer"
	}

	// void (no pointer) -> empty
	if ctype == "void" {
		return ""
	}

	// int -> int32
	if ctype == "int" {
		return "int32"
	}

	// size_t -> uintptr
	if ctype == "size_t" {
		return "uintptr"
	}

	// uintptr -> uintptr
	if ctype == "uintptr" {
		return "uintptr"
	}

	// double -> float64
	if ctype == "double" {
		return "float64"
	}

	// const char* -> *byte
	if ctype == "const char*" || ctype == "const char *" {
		return "*byte"
	}

	// char* -> *byte
	if ctype == "char*" || ctype == "char *" {
		return "*byte"
	}

	// Any pointer type -> unsafe.Pointer
	if isPtr {
		return "unsafe.Pointer"
	}

	// cef_string_t variants -> treat as uintptr (opaque)
	if strings.HasPrefix(ctype, "cef_string") {
		return "uintptr"
	}

	// Other cef_*_t types (non-pointer, likely enums or structs) -> uintptr
	if strings.HasPrefix(ctype, "cef_") {
		return "uintptr"
	}

	// Unknown non-pointer
	return "uintptr"
}
