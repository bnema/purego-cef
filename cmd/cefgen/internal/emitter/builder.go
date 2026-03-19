package emitter

import (
	"strings"

	"github.com/bnema/purego-cef/cmd/cefgen/internal/model"
)

// BuildPublicFileData converts a parsed header and type registry into the
// view model used by the public API templates.
func BuildPublicFileData(header *model.Header, registry *TypeRegistry) *PublicFileData {
	data := &PublicFileData{
		PackageName: "cef",
	}

	for i := range header.Structs {
		s := &header.Structs[i]
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
		data.FreeFunctions = append(data.FreeFunctions, buildFreeFunc(&header.Functions[i], registry))
	}

	return data
}

func buildInterface(s *model.Struct, registry *TypeRegistry) InterfaceData {
	iface := InterfaceData{
		Name:      s.InterfaceName,
		Doc:       s.Doc,
		Kind:      s.Kind,
		RawGoName: s.GoName,
	}

	for _, f := range s.Fields {
		if !f.IsFunction {
			continue
		}
		if strings.EqualFold(f.CName, "base") {
			continue
		}

		m := buildMethod(f, registry)
		iface.Methods = append(iface.Methods, m)
	}

	return iface
}

func buildMethod(f model.Field, registry *TypeRegistry) MethodData {
	m := MethodData{
		Name:         model.PublicName(f.CName),
		Doc:          f.Doc,
		RawFieldName: f.GoName,
	}

	// Build params, skipping the first "self" parameter.
	for i, p := range f.Params {
		if i == 0 {
			continue // skip self
		}
		m.Params = append(m.Params, ParamData{
			Name:       paramName(p),
			PublicType: registry.ResolvePublicType(p.CType),
			CType:      p.CType,
		})
	}

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
			PublicType: pubType,
			CType:      ret,
			IsBool:     isBool,
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
	}

	return ed
}

func buildFreeFunc(fn *model.Function, registry *TypeRegistry) FreeFuncData {
	ff := FreeFuncData{
		Name:      model.PublicName(fn.CName),
		Doc:       fn.Doc,
		RawGoName: fn.GoName,
	}

	for _, p := range fn.Params {
		ff.Params = append(ff.Params, ParamData{
			Name:       paramName(p),
			PublicType: registry.ResolvePublicType(p.CType),
			CType:      p.CType,
		})
	}

	ret := strings.TrimSpace(fn.ReturnCType)
	if ret == "" || ret == "void" {
		ff.Return = ReturnData{IsVoid: true, CType: ret}
	} else {
		ff.Return = ReturnData{
			PublicType: registry.ResolvePublicType(ret),
			CType:      ret,
		}
	}

	return ff
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
