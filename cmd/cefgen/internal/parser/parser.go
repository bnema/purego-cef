package parser

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	// Extract doc comments from raw source BEFORE stripping comments.
	rawSource := string(data)
	docIdx := buildDocIndex(rawSource)

	clean := stripComments(data)
	// Join continuation lines so multi-line declarations become single lines.
	clean = joinLines(clean)

	out := &model.Header{Path: path, Package: "capi"}
	for _, match := range structRE.FindAllSubmatch(clean, -1) {
		st, err := parseStruct(string(match[2]), string(match[3]))
		if err != nil {
			return nil, fmt.Errorf("%s: struct %s: %w", path, string(match[3]), err)
		}
		// Populate doc, kind, and interface name from raw source.
		populateStructDoc(&st, rawSource, docIdx)
		out.Structs = append(out.Structs, st)
	}
	for _, match := range funcRE.FindAllSubmatch(clean, -1) {
		fn := parseFunction(string(match[2]), string(match[1]), string(match[3]))
		// Populate function doc from raw source.
		populateFunctionDoc(&fn, rawSource, docIdx)
		out.Functions = append(out.Functions, fn)
	}
	for _, match := range enumRE.FindAllSubmatch(clean, -1) {
		out.Enums = append(out.Enums, parseEnum(string(match[2]), string(match[1])))
	}
	return out, nil
}

// populateStructDoc fills in Doc, Kind, and InterfaceName on a parsed struct
// by looking up the doc block in the raw source.
func populateStructDoc(st *model.Struct, rawSource string, docIdx *docIndex) {
	// Find the typedef line for this struct in the raw source.
	internalName := "_" + st.CName
	needle := "typedef struct " + internalName
	line := findLineOf(rawSource, needle)
	if line >= 0 {
		if db := docIdx.forLine(line); db != nil {
			st.Doc = cleanDoc(db.lines)
			kind := classifyKind(db.lines)
			if kind != "" {
				st.Kind = kind
			}
		}
	}
	// Default kind for structs without allocation comment.
	if st.Kind == "" {
		st.Kind = "data"
	}
	// Set InterfaceName from C name.
	st.InterfaceName = model.PublicName(st.CName)

	// Populate field docs from raw source.
	for i := range st.Fields {
		doc := findFieldDoc(rawSource, st.CName, st.Fields[i].CName)
		if doc != "" {
			st.Fields[i].Doc = doc
		}
	}
}

// populateFunctionDoc fills in Doc on a parsed function by looking up the
// doc block in the raw source.
func populateFunctionDoc(fn *model.Function, rawSource string, docIdx *docIndex) {
	// Find the CEF_EXPORT line for this function.
	needle := fn.CName + "("
	line := findLineOf(rawSource, needle)
	if line >= 0 {
		if db := docIdx.forLine(line); db != nil {
			fn.Doc = cleanDoc(db.lines)
		}
	}
}

// defineRE matches simple #define NAME VALUE lines (integer or expression).
var defineRE = regexp.MustCompile(`^#define\s+(\w+)\s+(\d+)\s*$`)

// stripComments removes single-line comments, block comments, preprocessor directives, and blank lines.
// It first collects simple #define constants and expands them inline.
func stripComments(data []byte) []byte {
	// First strip block comments /* ... */
	for {
		start := bytes.Index(data, []byte("/*"))
		if start < 0 {
			break
		}
		end := bytes.Index(data[start+2:], []byte("*/"))
		if end < 0 {
			break
		}
		data = append(data[:start], data[start+2+end+2:]...)
	}

	// Collect simple #define NAME INTEGER constants.
	defines := map[string]string{}
	for line := range bytes.SplitSeq(data, []byte("\n")) {
		if m := defineRE.FindSubmatch(bytes.TrimSpace(line)); m != nil {
			defines[string(m[1])] = string(m[2])
		}
	}

	// Expand #define constants in the source.
	for name, val := range defines {
		data = bytes.ReplaceAll(data, []byte(name), []byte(val))
	}

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

func parseStruct(body, name string) (model.Struct, error) {
	result := model.Struct{CName: name, GoName: goName(name)}
	for raw := range strings.SplitSeq(body, ";") {
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
		starCount := 0
		for {
			trimmed, ok := strings.CutPrefix(name, "*")
			if !ok {
				break
			}
			starCount++
			name = trimmed
		}
		ctype := strings.TrimSpace(strings.Join(fields[:len(fields)-1], " ")) + strings.Repeat("*", starCount)
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
	// Track values for auto-increment and symbolic resolution.
	nextVal := 0
	nameToVal := map[string]int{}
	seen := map[string]bool{}

	for line := range strings.SplitSeq(body, ",") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		cname := strings.TrimSpace(parts[0])
		if cname == "" {
			continue
		}

		// Skip duplicate enum values (from preprocessor conditionals).
		if seen[cname] {
			continue
		}
		seen[cname] = true

		val := ""
		if len(parts) == 2 {
			val = strings.TrimSpace(parts[1])
		}

		if val == "" {
			// Auto-increment from previous value.
			val = fmt.Sprintf("%d", nextVal)
			nameToVal[cname] = nextVal
			nextVal++
		} else {
			// Try to parse as integer.
			if n, err := strconv.Atoi(val); err == nil {
				nameToVal[cname] = n
				nextVal = n + 1
			} else if resolved, ok := nameToVal[val]; ok {
				// Symbolic reference to another enum value in same enum.
				val = fmt.Sprintf("%d", resolved)
				nameToVal[cname] = resolved
				nextVal = resolved + 1
			} else {
				// Resolve well-known C macros.
				val = resolveCMacros(val)
				// Complex expression (bit shifts, etc.) — keep as-is,
				// don't update auto-increment (next bare value will be wrong
				// but this is rare in practice).
				nameToVal[cname] = nextVal
				nextVal++
			}
		}

		result.Values = append(result.Values, model.EnumValue{
			CName:  cname,
			GoName: goName(cname),
			Value:  val,
		})
	}
	return result
}

// resolveCMacros replaces well-known C macros with Go equivalents.
func resolveCMacros(val string) string {
	replacements := map[string]string{
		"UINT_MAX": "0xFFFFFFFF",
	}
	for macro, replacement := range replacements {
		val = strings.ReplaceAll(val, macro, replacement)
	}
	return val
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
	ctype = strings.Join(strings.Fields(ctype), " ")

	isPtr := strings.Contains(ctype, "*")

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

	// float -> float32
	if ctype == "float" {
		return "float32"
	}

	// stdint types
	if ctype == "uint32_t" {
		return "uint32"
	}
	if ctype == "uint64_t" {
		return "uint64"
	}
	if ctype == "int64_t" {
		return "int64"
	}
	if ctype == "int32_t" {
		return "int32"
	}
	if ctype == "uint16_t" {
		return "uint16"
	}
	if ctype == "int16_t" {
		return "int16"
	}
	if ctype == "uint8_t" {
		return "uint8"
	}
	if ctype == "int8_t" {
		return "int8"
	}
	if ctype == "char16_t" {
		return "uint16"
	}

	// unsigned long -> uint64 (linux cef_window_handle_t etc.)
	if ctype == "unsigned long" {
		return "uint64"
	}

	// const char* -> *byte
	if ctype == "const char*" || ctype == "const char *" {
		return "*byte"
	}

	// char* -> *byte
	if ctype == "char*" || ctype == "char *" {
		return "*byte"
	}

	// char** -> uintptr
	if ctype == "char**" || ctype == "char **" {
		return "uintptr"
	}

	// Any pointer type -> unsafe.Pointer
	if isPtr {
		return "unsafe.Pointer"
	}

	// Opaque string handle types -> uintptr
	switch ctype {
	case "cef_string_userfree_t", "cef_string_userfree_utf8_t",
		"cef_string_userfree_utf16_t", "cef_string_userfree_wide_t",
		"cef_string_list_t", "cef_string_map_t", "cef_string_multimap_t":
		return "uintptr"
	}

	// Non-pointer cef_*_t types: use the Go struct/type name so embedded
	// structs expand to their full layout (critical for cef_base_ref_counted_t
	// and cef_string_t).
	if strings.HasPrefix(ctype, "cef_") && strings.HasSuffix(ctype, "_t") {
		return goName(ctype)
	}

	// Unknown non-pointer
	return "uintptr"
}
