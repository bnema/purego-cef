package emitter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

// rawConstructorTypes lists handler types whose generated constructor should
// be unexported (newRawXxx instead of NewXxx) because they have hand-written
// public constructors in bridge.go that wrap them with safe types.
var rawConstructorTypes = map[string]bool{
	"Client":          true,
	"LifeSpanHandler": true,
}

// skipPublicTypes lists type/function names (after PublicName conversion) that
// are hand-written in init.go or support.go and must not be generated.
var skipPublicTypes = map[string]bool{
	"Settings":          true,
	"MainArgs":          true,
	"Shutdown":          true,
	"DoMessageLoopWork": true,
	"Initialize":        true,
	"ExecuteProcess":    true,
	// Functions with struct-by-value params/returns need hand-written wrappers.
	"V8ValueCreateDate":                   true,
	"DisplayConvertScreenPointToPixels":   true,
	"DisplayConvertScreenPointFromPixels": true,
	"DisplayConvertScreenRectToPixels":    true,
	"DisplayConvertScreenRectFromPixels":  true,
	// Handlers with unsafe out-params that need hand-written safe signatures.
	"LifeSpanHandler": true,
	"AudioHandler":    true,
}

// ParamOverride specifies the public type and marshal kind for a specific
// handler callback parameter, overriding the default type resolution.
type ParamOverride struct {
	PublicType    string
	MarshalKind   string
	UnmarshalExpr string
}

// paramOverrides maps "structCName.fieldCName.paramCName" keys to overrides
// for specific handler callback params that need safe public types instead of
// unsafe.Pointer.
var paramOverrides = map[string]ParamOverride{
	// Problem 5: GetScreenPoint out-params
	"cef_render_handler_t.get_screen_point.screenX": {PublicType: "*int32", MarshalKind: "dataStruct"},
	"cef_render_handler_t.get_screen_point.screenY": {PublicType: "*int32", MarshalKind: "dataStruct"},
	// Problem 6: OnPaint pixel buffer
	"cef_render_handler_t.on_paint.buffer": {
		PublicType:    "[]byte",
		MarshalKind:   "pixelBuffer",
		UnmarshalExpr: "int({{width}})*int({{height}})*4",
	},
	// Problem 8: GetResourceRequestHandler bool flag
	"cef_request_handler_t.get_resource_request_handler.disable_default_handling": {PublicType: "*int32", MarshalKind: "dataStruct"},
	// Problem 9: OnSelectClientCertificate object array
	"cef_request_handler_t.on_select_client_certificate.certificates": {PublicType: "[]X509Certificate", MarshalKind: "objectSlice"},
}

// BuildPublicFileData converts a parsed header and type registry into the
// view model used by the public API templates. Types in skipPublicTypes are
// excluded (they have handwritten implementations).
func BuildPublicFileData(header *model.Header, registry *TypeRegistry) *PublicFileData {
	return buildFileData(header, registry, true)
}

// BuildPortFileData converts a parsed header into a view model for port
// templates. Unlike BuildPublicFileData, it does NOT skip any types — ports
// need the full interface surface so cross-package references resolve.
func BuildPortFileData(header *model.Header, registry *TypeRegistry) *PublicFileData {
	return buildFileData(header, registry, false)
}

func buildFileData(header *model.Header, registry *TypeRegistry, applySkip bool) *PublicFileData {
	stem := strings.TrimPrefix(header.RegisterName, "Register")
	data := &PublicFileData{
		PackageName: "cef",
		HeaderStem:  stem,
	}

	for i := range header.Structs {
		s := &header.Structs[i]
		pubName := model.PublicName(s.CName)
		if applySkip && skipPublicTypes[pubName] {
			continue
		}
		switch s.Kind {
		case "handler", "object":
			data.Interfaces = append(data.Interfaces, buildInterface(s, registry))
		case "data":
			data.DataStructs = append(data.DataStructs, buildDataStruct(s, registry))
		}
	}

	for i := range header.Enums {
		data.Enums = append(data.Enums, buildEnum(&header.Enums[i]))
	}

	for i := range header.Functions {
		fn := &header.Functions[i]
		pubName := model.PublicName(fn.CName)
		if applySkip && skipPublicTypes[pubName] {
			continue
		}
		data.FreeFunctions = append(data.FreeFunctions, buildFreeFunc(fn, registry))
	}

	return data
}

func buildInterface(s *model.Struct, registry *TypeRegistry) InterfaceData {
	// Detect scoped types by checking the base field type.
	isScoped := false
	for _, f := range s.Fields {
		if strings.EqualFold(f.CName, "base") && strings.Contains(f.CType, "cef_base_scoped_t") {
			isScoped = true
			break
		}
	}

	iface := InterfaceData{
		Name:      s.InterfaceName,
		Doc:       s.Doc,
		Kind:      s.Kind,
		IsScoped:  isScoped,
		RawGoName: s.GoName,
	}

	for _, f := range s.Fields {
		if !f.IsFunction {
			continue
		}
		if strings.EqualFold(f.CName, "base") {
			continue
		}

		m := buildMethod(s.CName, f, registry)
		iface.Methods = append(iface.Methods, m)
	}

	return iface
}

func buildMethod(structCName string, f model.Field, registry *TypeRegistry) MethodData {
	name := model.PublicName(f.CName)
	if renamed, ok := methodRenames[name]; ok {
		name = renamed
	}
	m := MethodData{
		Name:         name,
		Doc:          f.Doc,
		RawFieldName: f.GoName,
	}

	// Build params, skipping the first "self" parameter.
	for i, p := range f.Params {
		if i == 0 {
			continue // skip self
		}
		pd := ParamData{
			Name:        paramName(p),
			PublicType:  registry.ResolvePublicType(p.CType),
			CType:       p.CType,
			MarshalKind: classifyParamType(p.CType, registry),
		}
		m.Params = append(m.Params, pd)
	}

	// Apply param overrides before copying to RawParams and merging.
	for i := range m.Params {
		key := structCName + "." + f.CName + "." + f.Params[i+1].CName // +1 to skip self
		if ov, ok := paramOverrides[key]; ok {
			m.Params[i].PublicType = ov.PublicType
			m.Params[i].MarshalKind = ov.MarshalKind
		}
	}

	// Resolve UnmarshalExpr placeholders to raw arg names.
	for i := range m.Params {
		key := structCName + "." + f.CName + "." + f.Params[i+1].CName
		if ov, ok := paramOverrides[key]; ok && ov.UnmarshalExpr != "" {
			expr := ov.UnmarshalExpr
			for j, rp := range m.Params {
				placeholder := "{{" + rp.Name + "}}"
				argName := fmt.Sprintf("arg%d", j)
				expr = strings.ReplaceAll(expr, placeholder, argName)
			}
			m.Params[i].UnmarshalExtra = expr
		}
	}

	// Save original params for raw callback signatures, then merge count+pointer pairs.
	m.RawParams = make([]ParamData, len(m.Params))
	copy(m.RawParams, m.Params)
	m.Params = mergeCountPointerParams(m.Params, registry)

	// Build return type.
	ret := strings.TrimSpace(f.ReturnCType)
	if ret == "" || ret == "void" {
		m.Return = ReturnData{IsVoid: true, CType: ret}
	} else {
		isBool := IsBoolReturn(f)
		pubType := registry.ResolvePublicType(ret)
		if isBool {
			pubType = "bool"
		}
		m.Return = ReturnData{
			PublicType:   pubType,
			CType:        ret,
			IsBool:       isBool,
			IsEnum:       registry.IsEnumType(ret),
			IsString:     registry.IsStringType(ret),
			IsInterface:  registry.IsInterfaceType(ret),
			IsNumeric:    isNumericType(pubType),
			IsPointer:    pubType == "unsafe.Pointer",
			IsHandler:    registry.IsHandlerType(ret),
			IsDataStruct: registry.IsDataStructType(ret),
		}
	}

	// Check if this is a getter callback.
	if registry.IsGetterCallback(f) {
		m.IsGetter = true
		m.GetterInterface = registry.ResolvePublicType(f.ReturnCType)
	}

	return m
}

func buildDataStruct(s *model.Struct, registry *TypeRegistry) DataStructData {
	ds := DataStructData{
		Name:      model.PublicName(s.CName),
		Doc:       s.Doc,
		RawGoName: s.GoName,
	}

	for _, f := range s.Fields {
		if strings.EqualFold(f.CName, "base") {
			continue
		}
		ds.Fields = append(ds.Fields, DataFieldData{
			Name:       model.PublicName(f.CName),
			PublicType: registry.ResolvePublicType(f.CType),
			Doc:        f.Doc,
		})
	}

	// Sized data structs need a New*() constructor that sets Size = unsafe.Sizeof(v).
	// Base types (cef_base_ref_counted_t etc.) also have a size field but their
	// Size is set by initRefCount, so exclude structs with function pointers.
	if !s.HasFunctionFields() {
		for _, f := range s.Fields {
			if f.CName == "size" && f.CType == "size_t" {
				ds.HasSizeField = true
				break
			}
		}
	}

	return ds
}

func buildEnum(e *model.Enum) EnumData {
	ed := EnumData{
		Name:      model.PublicName(e.CName),
		RawGoName: e.GoName,
	}

	// Determine the common prefix to strip from enum values.
	// e.g., for cef_state_t values like STATE_DEFAULT, STATE_ENABLED, etc.
	prefix := enumValuePrefix(e)

	for _, v := range e.Values {
		name := cleanEnumValueName(v.CName, prefix)
		ed.Values = append(ed.Values, EnumValueData{
			Name:  name,
			Value: v.Value,
		})
		// Detect unsigned values (hex > 0x7FFFFFFF or large decimal).
		if isUnsignedValue(v.Value) {
			ed.Unsigned = true
		}
	}

	return ed
}

// isUnsignedValue returns true if the value string represents a number
// that exceeds the int32 range.
func isUnsignedValue(val string) bool {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		// Parse as uint64 and check if > MaxInt32.
		n, err := strconv.ParseUint(val[2:], 16, 64)
		if err == nil && n > 0x7FFFFFFF {
			return true
		}
	}
	return false
}

func buildFreeFunc(fn *model.Function, registry *TypeRegistry) FreeFuncData {
	ff := FreeFuncData{
		Name:      model.PublicName(fn.CName),
		Doc:       fn.Doc,
		RawGoName: fn.GoName,
	}

	for _, p := range fn.Params {
		ff.Params = append(ff.Params, ParamData{
			Name:        paramName(p),
			PublicType:  registry.ResolvePublicType(p.CType),
			CType:       p.CType,
			MarshalKind: classifyParamType(p.CType, registry),
			RawGoType:   p.GoType,
		})
	}

	ret := strings.TrimSpace(fn.ReturnCType)
	if ret == "" || ret == "void" {
		ff.Return = ReturnData{IsVoid: true, CType: ret}
	} else {
		pubType := registry.ResolvePublicType(ret)
		ff.Return = ReturnData{
			PublicType:  pubType,
			CType:       ret,
			IsInterface: registry.IsInterfaceType(ret),
			IsString:    registry.IsStringType(ret),
			IsEnum:      registry.IsEnumType(ret),
			IsNumeric:   isNumericType(pubType),
			IsPointer:   pubType == "unsafe.Pointer",
			IsBool:      false, // free functions don't use bool heuristic
		}
	}

	return ff
}

// methodRenames maps method names that conflict with Go standard interfaces
// (like io.Seeker's Seek) to alternative names to avoid go vet warnings.
var methodRenames = map[string]string{
	"Seek": "SeekOffset",
}

// goKeywords are Go reserved words that cannot be used as identifiers.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true,
	"for": true, "func": true, "go": true, "goto": true, "if": true,
	"import": true, "interface": true, "map": true, "package": true,
	"range": true, "return": true, "select": true, "struct": true,
	"switch": true, "type": true, "var": true,
}

// paramName converts a C parameter name to a Go-idiomatic name.
func paramName(p model.Param) string {
	name := p.GoName
	if name == "" {
		name = model.PublicName(p.CName)
	}
	// Lowercase first letter for parameter names.
	if len(name) > 0 {
		name = strings.ToLower(name[:1]) + name[1:]
	}
	// Escape Go keywords by appending an underscore.
	if goKeywords[name] {
		name = name + "_"
	}
	return name
}

// enumValuePrefix finds the common prefix of all enum value C names.
// For CEF enums, values typically look like CEF_STATE_DEFAULT, CEF_STATE_ENABLED.
// We strip "CEF_" and the type-specific prefix.
func enumValuePrefix(e *model.Enum) string {
	if len(e.Values) == 0 {
		return ""
	}

	// Strip "cef_" prefix and "_t" suffix to get the enum stem, e.g. "state".
	stem := strings.TrimLeft(e.CName, "_")
	stem = strings.TrimPrefix(stem, "cef_")
	stem = strings.TrimSuffix(stem, "_t")

	// The enum value prefix is usually "CEF_" + uppercase stem + "_".
	// e.g., stem = "state" → prefix = "CEF_STATE_"
	// But also handle: "CEFV8_" style prefixes.
	prefix := "CEF_" + strings.ToUpper(stem) + "_"

	// Verify this prefix works for at least one value.
	for _, v := range e.Values {
		if strings.HasPrefix(v.CName, prefix) {
			return prefix
		}
	}

	// Fallback: just strip "CEF_".
	return "CEF_"
}

// classifyParamType returns a MarshalKind string for a given C type.
func classifyParamType(ctype string, registry *TypeRegistry) string {
	ct := normalizeConst(strings.TrimSpace(ctype))
	switch ct {
	case "const cef_string_t*", "const char*", "char*":
		return "string"
	case "cef_string_userfree_t":
		return "userfreeString"
	case "void*", "const void*":
		return "pointer"
	case "cef_string_t*", "cef_string_list_t", "cef_string_map_t", "cef_string_multimap_t":
		return "numeric" // opaque handles, passed as uintptr
	}
	if registry.IsEnumType(ct) {
		return "enum"
	}
	if registry.IsInterfaceType(ct) {
		return "interface"
	}
	if registry.IsDataStructType(ct) {
		return "dataStruct"
	}
	pub := registry.ResolvePublicType(ct)
	if pub == "unsafe.Pointer" {
		return "pointer"
	}
	if strings.HasPrefix(pub, "*") {
		return "dataStruct" // pointer to struct/int/etc, cast as (*Type)(unsafe.Pointer(...))
	}
	if isNumericType(pub) {
		return "numeric"
	}
	return "numeric" // default fallback
}

// isNumericType returns true if the Go type is a numeric primitive.
func isNumericType(goType string) bool {
	switch goType {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "float32", "float64":
		return true
	}
	return false
}

// cleanEnumValueName strips the common prefix and PascalCases the result.
func cleanEnumValueName(cName, prefix string) string {
	name := cName
	if strings.HasPrefix(name, prefix) {
		name = name[len(prefix):]
	} else if strings.HasPrefix(name, "CEF_") {
		name = name[4:]
	}

	// Convert UPPER_SNAKE to PascalCase.
	parts := strings.Split(strings.ToLower(name), "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	result := b.String()
	if result == "" {
		return cName // fallback to original
	}
	return result
}

// mergeCountPointerParams detects count+pointer param pairs in handler callbacks
// and merges them into a single []ElemType slice param.
//
// Pattern: param[i] has name ending in "count" with numeric type, and param[i+1]
// has a name matching the prefix (e.g., "dirtyrectscount" + "dirtyrects").
// The pointer param's type must resolve to a data struct pointer ("*Rect" etc.).
func mergeCountPointerParams(params []ParamData, registry *TypeRegistry) []ParamData {
	if len(params) < 2 {
		return params
	}

	var merged []ParamData
	skip := false
	for i := 0; i < len(params); i++ {
		if skip {
			skip = false
			continue
		}
		if i+1 < len(params) {
			countP := params[i]
			ptrP := params[i+1]

			// Check: count param is numeric, name ends with "count",
			// and the pointer param's name matches the prefix.
			countName := strings.ToLower(countP.Name)
			ptrName := strings.ToLower(ptrP.Name)
			if strings.HasSuffix(countName, "count") &&
				isIntLikeType(countP.PublicType) &&
				strings.HasPrefix(countName, ptrName) &&
				(ptrP.MarshalKind == "dataStruct" || ptrP.MarshalKind == "objectSlice") {

				// Extract element type: "*Rect" → "Rect", "[]X509Certificate" → "X509Certificate"
				elemType := strings.TrimPrefix(ptrP.PublicType, "*")
				elemType = strings.TrimPrefix(elemType, "[]")

				// Preserve objectSlice kind; default to slice for dataStruct.
				marshalKind := "slice"
				if ptrP.MarshalKind == "objectSlice" {
					marshalKind = "objectSlice"
				}

				sliceParam := ParamData{
					Name:          ptrP.Name,
					PublicType:    "[]" + elemType,
					CType:         ptrP.CType,
					MarshalKind:   marshalKind,
					SliceElemType: elemType,
					// SliceCountArg/SlicePtrArg set by template using RawParams indices
				}
				merged = append(merged, sliceParam)
				skip = true
				continue
			}
		}
		merged = append(merged, params[i])
	}
	return merged
}

// isIntLikeType returns true if the Go type is an integer type suitable as an array count.
func isIntLikeType(goType string) bool {
	switch goType {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return true
	}
	return false
}
